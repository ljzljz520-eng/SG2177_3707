package main

import (
	"context"
	"fmt"
	"io"

	"equipmentlending/internal/reporting"
	"equipmentlending/internal/service"
)

func runCreate(ctx context.Context, business *service.Service, writer io.Writer, request service.CreateRequest) error {
	record, err := business.Create(ctx, request)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, "created %s\n", record.Summary())
	return err
}

func runQuery(ctx context.Context, business *service.Service, writer io.Writer, borrower string) error {
	records, err := business.Query(ctx, borrower, "")
	if err != nil {
		return err
	}
	return reporting.Build(records, "Equipment Lending Records").Write(writer)
}

func runSort(ctx context.Context, business *service.Service, writer io.Writer, mode string) error {
	records, err := business.Sort(ctx, mode)
	if err != nil {
		return err
	}
	return reporting.Build(records, "Sorted Equipment").Write(writer)
}
