package main

import (
	"github.com/mizk/zzdm"
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var (
	secret   = false
	force    = false
	version  = false
	advice   = false
	password = ""
	input    = ""
	output   = ""
)

const (
	ROOT       = iota
	ENCRYPTION
	DECRYPTION
)

func main() {

	command := &cobra.Command{Use: "zzdm",
		Short: "zzdm is a file encryption/decryption tool with aes crypt",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				fmt.Printf("zzdm-%s-%s\n", zzdm.Version, zzdm.SKU)
			} else {
				cmd.Println(cmd.UsageString())
			}
		}}
	parseFlag(command, ROOT)
	encrypt := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a file",
		Long:  "zzdm encrypt [-s | --secret] [-a | --advice] [-f | --force] (-i | --input $input) [-o | --output $output] (-p | --password $password)",
		Run: func(cmd *cobra.Command, args []string) {
			if !zzdm.Exist(input) {
				fmt.Println("input file is missing")
				os.Exit(-1)
				return
			}
			if !zzdm.IsDir(output) {
				output = ""
			}
			if len(password) == 0 {
				fmt.Println("password is required")
				os.Exit(-2)
				return
			}
			if advice {
				checkPassword(password)
			}
			err := zzdm.Encrypt(input, output, password, secret, force)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		},
	}
	parseFlag(encrypt, ENCRYPTION)
	command.AddCommand(encrypt)

	decrypt := &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt a file",
		Long:  "zzdm decrypt [--force] (-i|--input $input) [-o|--output $output] (-p|--password $password)",
		Run: func(cmd *cobra.Command, args []string) {
			if !zzdm.Exist(input) {
				fmt.Println("input file is missing")
				os.Exit(-1)
				return
			}
			if !zzdm.IsDir(output) {
				output = ""
			}
			if len(password) == 0 {
				fmt.Print("password is required")
				os.Exit(-2)
				return
			}
			err := zzdm.Decrypt(input, output, password, force)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		},
	}
	parseFlag(decrypt, DECRYPTION)
	command.AddCommand(decrypt)
	err := command.Execute()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func checkPassword(password string) {
	level := zzdm.PasswordLevel(password)
	if level == 0 {
		fmt.Println("[PA]:Perfect")
	} else if level == -1 {
		fmt.Println("[PA]:Password should contain uppercase letters")
	} else if level == -2 {
		fmt.Println("[PA]:Password should contain lowercase letters")
	} else if level == -3 {
		fmt.Println("[PA]:Password should contain numbers")
	} else if level == -4 {
		fmt.Println("[PA]:Password should be at least 8 characters in length")
	} else if level == -5 {
		fmt.Println("[PA]:Password should contain at least 1 special character(~`!@#$%^&*()_+-=[]{}|\\<,>.?/;:\"')")
	}

}

func parseFlag(command *cobra.Command, classify int) {
	if classify == ROOT {
		command.PersistentFlags().BoolVarP(&version, "version", "v", false, "display version info")
	} else {
		command.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
		command.PersistentFlags().StringVarP(&output, "output", "o", "", "output directory")
		command.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
		command.PersistentFlags().BoolVarP(&force, "force", "f", false, "force to overwrite an existing  file within the output directory")
		if classify == ENCRYPTION {
			command.PersistentFlags().BoolVarP(&advice, "advice", "a", false, "show password advice")
			command.PersistentFlags().BoolVarP(&secret, "secret", "s", false, "random output file name")
		}

	}

}
