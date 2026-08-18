package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func runCLI(args []string) int {
	application := newApplication()
	if len(args) == 0 {
		fmt.Println("calibration-vault requires a command")
		return 2
	}
	switch args[0] {
	case "describe":
		return describeCLI()
	case "check-intake":
		return checkIntakeCLI(application)
	case "check-review":
		return checkReviewCLI(application)
	case "check-release":
		return checkReleaseCLI(application)
	case "check-query":
		return checkQueryCLI(application)
	case "check-batch":
		return checkBatchCLI(application)
	default:
		fmt.Printf("unknown command: %s\n", args[0])
		return 2
	}
}

func newApplication() *service.Service {
	return service.New(store.NewMemory(), service.SystemClock{})
}

func describeCLI() int {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("input is required")
		return 2
	}
	operator := strings.TrimSpace(scanner.Text())
	if operator == "" {
		fmt.Println("operator is required")
		return 2
	}
	fmt.Printf("ready for %s\n", operator)
	return 0
}

func runtimeIntake() domain.IntakeInput {
	return domain.IntakeInput{
		SampleID: "runtime-sample",
		Source:   "line-a",
		Priority: domain.PriorityPriority,
		Tags:     []string{"calibration", "runtime"},
	}
}

func runtimeReview() domain.ReviewInput {
	return domain.ReviewInput{
		Reviewer: "runtime-reviewer",
		Decision: domain.DecisionApprove,
		Note:     "metadata confirmed",
	}
}

func runtimeRelease() domain.ReleaseInput {
	return domain.ReleaseInput{
		Operator:    "runtime-release",
		Destination: "lab-west",
	}
}

func checkIntakeCLI(application *service.Service) int {
	sample, err := application.Intake(context.Background(), runtimeIntake())
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("created %s\n", sample.Status)
	return 0
}

func checkReviewCLI(application *service.Service) int {
	sample, err := application.Intake(context.Background(), runtimeIntake())
	if err == nil {
		sample, err = application.Review(context.Background(), sample.SampleID, runtimeReview())
	}
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("reviewed %s\n", sample.Status)
	return 0
}

func checkReleaseCLI(application *service.Service) int {
	result, err := application.IntakeApproveRelease(
		context.Background(),
		runtimeIntake(),
		runtimeReview(),
		runtimeRelease(),
	)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("released %s\n", result.Sample.Release.Destination)
	return 0
}

func checkQueryCLI(application *service.Service) int {
	result, err := application.IntakeApproveRelease(
		context.Background(),
		runtimeIntake(),
		runtimeReview(),
		runtimeRelease(),
	)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("timeline events=%d\n", len(result.Timeline))
	return 0
}

func checkBatchCLI(application *service.Service) int {
	input := domain.BatchIntakeInput{
		BatchReference: "runtime-batch",
		Items: []domain.IntakeInput{
			{SampleID: "batch-one", Source: "line-a", Priority: domain.PriorityRoutine},
			{SampleID: "batch-two", Source: "line-b", Priority: domain.PriorityUrgent},
		},
	}
	result, err := application.BatchIntake(context.Background(), input)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("batch created count=%d\n", result.Count)
	return 0
}
