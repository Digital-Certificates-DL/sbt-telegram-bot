package service

import (
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/config"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/helpers"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/sbt"
)

func Start(cfg config.Config) error {
	//todo add description
	cfg.Log().Info("mint: 1\n" +
		"transfer: 2\n" +
		"burn: 3\n" +
		"add admin: 4\n" +
		"delete admin: 5\n")
	command, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to start")
	}
	caller, err := sbt.NewCaller(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create caller")
	}
	switch {
	case command == "1":
		err = caller.Mint()
		if err != nil {
			return errors.Wrap(err, "failed to mint")
		}
	case command == "2":
		err = caller.Transfer()
		if err != nil {
			return errors.Wrap(err, "failed to transfer")
		}
	case command == "3":
		err = caller.Burn()
		if err != nil {
			return errors.Wrap(err, "failed to burn")
		}
	case command == "4":
		err = caller.NewAdmin()
		if err != nil {
			return errors.Wrap(err, "failed to add new admin")
		}
	case command == "5":
		err = caller.DeleteAdmin()
		if err != nil {
			return errors.Wrap(err, "failed to delete admin")
		}
	case command == "6":
		err = caller.OwnerOf()
		if err != nil {
			return errors.Wrap(err, "failed to get owner of token")
		}
	case command == "7":
		err = caller.Name()
		if err != nil {
			return errors.Wrap(err, "failed to get name")
		}
	}

	return nil
}
