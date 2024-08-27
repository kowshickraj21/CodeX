package controllers

import (
	"fmt"
	"io/ioutil"
	"main/models"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func JavaExecuter(ctx *gin.Context) {
	var req models.Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Println("Received code:", req.Code)

	// Use fixed file names for Java source file and input file
	sourceFileName := "Main.java"
	inputFileName := "input.txt"

	// Create or overwrite the source file in the current directory
	if err := ioutil.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to write code to file"})
		return
	}
	defer os.Remove(sourceFileName) // Ensure the source file is deleted after execution

	// If there's input data, write it to a temporary input file
	if req.Input != "" {
		if err := ioutil.WriteFile(inputFileName, []byte(req.Input), 0644); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to write input to file"})
			return
		}
		defer os.Remove(inputFileName) // Ensure the input file is deleted after execution
	}

	fmt.Println("Temporary files created at:", sourceFileName, inputFileName)

	// Compile the Java code
	compileCmd := exec.Command("javac", sourceFileName)
	compileOutput, err := compileCmd.CombinedOutput()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Compilation failed", "output": string(compileOutput)})
		return
	}

	fmt.Println("Compilation completed")

	// Run the compiled Java code with input redirection if needed
	var runCmd *exec.Cmd
	if req.Input != "" {
		runCmd = exec.Command("java", "Main")
		runCmd.Stdin, _ = os.Open(inputFileName) // Redirect the input file to the Java process
	} else {
		runCmd = exec.Command("java", "Main")
	}
	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Execution failed", "output": string(runOutput)})
		return
	}
	defer os.Remove("Main.class") // Ensure the class file is deleted after execution

	ctx.JSON(200, gin.H{"output": string(runOutput)})
}