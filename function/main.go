package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

func main() {
	lambda.Start(handler)
}

// handler responsible to run the business logic for the scheduled
// cron job based on the cloudwatch events
func handler(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) error {
	users, err := fetchUsers()
	if err != nil {
		return err
	}
	maxTime := time.Now().AddDate(0, 0, 90)
	userToDeleteProfile, keysToDelete, err := checkCredentials(maxTime, users)
	if err != nil {
		return err
	}
	err = revoke(userToDeleteProfile, keysToDelete)
	if err != nil {
		return err
	}
	return nil
}

// fetchUsers retrieve the available users of the given config
func fetchUsers() (*iam.ListUsersOutput, error) {
	log.Info("Collecting users")
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return &iam.ListUsersOutput{}, errors.Wrap(err, "session.NewSession")
	}

	svc := iam.New(sess)
	result, err := svc.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeServiceFailureException:
				return &iam.ListUsersOutput{}, errors.Wrapf(aerr, "svc.ListUsers: %s", iam.ErrCodeServiceFailureException)
			default:
				return &iam.ListUsersOutput{}, errors.Wrap(aerr, "svc.ListUsers")
			}
		}
		return &iam.ListUsersOutput{}, errors.Wrapf(err, "svc.ListUsers")
	}
	return result, nil
}

// checkCredentials collect which account login credentials
// are unused for max time we have set when we ran lambda
func checkCredentials(maxTime time.Time, users *iam.ListUsersOutput) ([]string, []string, error) {
	log.Info("Checking Access Keys and Login Profile for users")
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return []string{}, []string{}, errors.Wrap(err, "session.NewSession")
	}
	created := time.Now()
	lastUsed := time.Now()

	var usersToDeleteProfile, accessKeysToDelete []string
	for _, u := range users.Users {
		// checks login profiles first
		if u.PasswordLastUsed != nil {
			lastUsed = *u.PasswordLastUsed
			diffDays := maxTime.Sub(lastUsed).Hours() / 24

			if diffDays > float64(maxTime.Day()) {
				log.Infof("User with ARN: %s, hasn't used credentials for %d", *u.Arn, int64(diffDays))
				usersToDeleteProfile = append(usersToDeleteProfile, *u.UserName)
			}
		}

		// checks access keys next
		svc := iam.New(sess)
		keys, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
			UserName: u.UserName,
			MaxItems: aws.Int64(2),
		})
		if err != nil {
			return []string{}, []string{}, errors.Wrap(err, "svc.ListAccessKeys")
		}

		for _, k := range keys.AccessKeyMetadata {
			if k.CreateDate != nil {
				created = *k.CreateDate
			}

			if k.AccessKeyId != nil {
				res, err := svc.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
					AccessKeyId: k.AccessKeyId,
				})
				if err != nil {
					return []string{}, []string{}, errors.Wrap(err, "svc.GetAccessKeyLastUsed")
				}
				if res != nil && res.AccessKeyLastUsed != nil && res.AccessKeyLastUsed.LastUsedDate != nil {
					lastUsed = *res.AccessKeyLastUsed.LastUsedDate
					continue
				}
				lastUsed = created
			}

			diffDays := maxTime.Sub(lastUsed).Hours() / 24
			if diffDays > float64(maxTime.Day()) {
				log.Infof("User with ARN: %s, hasn't used AccessKeyID: %s for %d", *u.Arn, *k.AccessKeyId, int64(diffDays))
				accessKeysToDelete = append(accessKeysToDelete, *k.AccessKeyId)
			}

		}
	}
	return usersToDeleteProfile, accessKeysToDelete, nil
}

// revoke the login profiles and access keys for the AWS users
func revoke(users, keys []string) error {
	log.Infof("Revoking login profiles for %d users", len(users))
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return errors.Wrap(err, "session.NewSession")
	}
	svc := iam.New(sess)
	for _, u := range users {
		_, err = svc.DeleteLoginProfile(&iam.DeleteLoginProfileInput{
			UserName: aws.String(u),
		})
		if err != nil {
			log.WithError(err).Error("Unable to delete login profile for user: %s", u)
		}
	}

	log.Infof("Revoking %d access keys ", len(keys))
	for _, k := range keys {
		_, err = svc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(k),
		})
		if err != nil {
			log.WithError(err).Error("Unable to delete key with ID: %s", k)
		}
	}
	return nil
}
