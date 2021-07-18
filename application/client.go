package application

import (
	"os"
	"os/exec"
	"io/ioutil"
	"fmt"
	"time"
)

const MEMO_EDITOR_INPUT_LEN = 2
const EDITOR_COMMAND_INPUT_POS = 1

const MEMO_INPUT_LEN = 3
const SAVE_COMMAND_INPUT_POS = 1
const MEMO_INPUT_POS = 2
const SAVE_COMMAND = "save"
const SAVE_BY_VIM_COMMAND = "vim"
const SAVE_BY_NVIM_COMMAND = "nvim"
const SAVE_BY_EMACS_COMMAND = "emacs"
const EDITOR_TEMP_FILE_DIR = ".flomo-tmp"
const VIM_TEMP_FILE_PREFIX = EDITOR_TEMP_FILE_DIR + "/vim-memo-tmp-"
const NVIM_TEMP_FILE_PREFIX = EDITOR_TEMP_FILE_DIR + "/nvim-memo-tmp-"
const EMACS_TEMP_FILE_PREFIX = EDITOR_TEMP_FILE_DIR + "/emacs-memo-tmp-"
var EDITOR_COMMAND_MAP = map[string]string{
    SAVE_BY_VIM_COMMAND: VIM_TEMP_FILE_PREFIX,
    SAVE_BY_NVIM_COMMAND: NVIM_TEMP_FILE_PREFIX,
    SAVE_BY_EMACS_COMMAND: EMACS_TEMP_FILE_PREFIX,
}
const CLEAR_TEMP_COMMAND = "clear"

const CONFIG_INPUT_LEN = 4
const SET_INPUT_POS = 1
const SET_COMMAND_INPUT_POS = 2
const COMMAND_CONTENT_INPUT_POS = 3
const SET_COMMAND = "set"
const SET_API_COMMAND = "api"

// Handle is the entry of the cli.
func Handle() {
	input := os.Args
	length := len(input)
	switch length {
	case MEMO_EDITOR_INPUT_LEN:
		editorCommand := input[EDITOR_COMMAND_INPUT_POS]
		if editorCommand == CLEAR_TEMP_COMMAND {
			removeTmpFile(ParseFilePath(EDITOR_TEMP_FILE_DIR), true)
			fmt.Println("All temporary files of the editors are cleared!")
			return
		}
		memo := execTextEditor(editorCommand)
		SendMemo(memo)
	case MEMO_INPUT_LEN:
		if input[SAVE_COMMAND_INPUT_POS] != SAVE_COMMAND {
			fmt.Println("Invalid command, use save to send a memo to flomo.")
			return
		}
		memo := input[MEMO_INPUT_POS]
		if len(memo) == 0 {
			fmt.Println("The memo you typed is blank, try to say some thing.")
			return
		}
		SendMemo(memo)
	case CONFIG_INPUT_LEN:
		if input[SET_INPUT_POS] != SET_COMMAND {
			fmt.Println("Invalid command, use set to configure the flomo-cli.")
			return
		}
		if input[SET_COMMAND_INPUT_POS] == SET_API_COMMAND {
			commandContent := input[COMMAND_CONTENT_INPUT_POS]
			if commandContent == "" {
				fmt.Println("The command content is blank, try [flomo set api].")
			}
			SaveConfig(FlomoConfig{Api: commandContent})
			return
		}
		fmt.Println("Nothing to set, try [flomo set api].")
	default:
		fmt.Println("Invalid input. Use set to configure the cli or use save to type a new memo.")
	}
}

// Exec the required editor and read editted data from the temp file.
// For security, only matched editor can be executed.
func execTextEditor(editorType string) string {
	filePrefix := ParseFilePath(EDITOR_COMMAND_MAP[editorType])
	if filePrefix == "" {
		panic("The editor you typed is not currently supported, try to use vim/nvim/emacs instead.")
	}
	filePath := filePrefix + time.Now().Local().String()

	err := os.MkdirAll(ParseFilePath(EDITOR_TEMP_FILE_DIR), os.ModePerm)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(editorType, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	tmpData, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	defer removeTmpFile(filePath, false)
	return string(tmpData)
}

// Remove or clear the temp editor files.
func removeTmpFile(filePath string, removeAll bool) {
	realFilePath := filePath
	var err error
	if removeAll {
		err = os.RemoveAll(realFilePath)
	} else {
		err = os.Remove(realFilePath)
	}
	if err != nil {
		panic(err)
	}
}
