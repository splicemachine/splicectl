package common

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestNoDatabaseName(t *testing.T) {
	tcmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			name := DatabaseName(cmd)
			if name != "" {
				t.Fatalf("there should have been no database name, but instead found: %s", name)
			}
		},
	}
	tcmd.SetArgs([]string{})
	if err := tcmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseName(t *testing.T) {
	for _, args := range [][]string{
		{"-d", "splicedb"},
		{"--database-name", "splicedb"},
		{"--workspace", "splicedb"},
		{"--database", "splicedb"},
		{"--database-name", "splicedb", "--workspace", "workspace", "--database", "database"},
		{"--database-name", "splicedb", "--workspace", "workspace"},
		{"--database-name", "splicedb", "--database", "database"},
		{"--workspace", "splicedb", "--database", "workspace"},
	} {
		testDatabaseNameWithArray(t, args)
	}
}

func testDatabaseNameWithArray(t *testing.T, args []string) {
	tcmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Visit(func(f *pflag.Flag) { fmt.Printf("%s: %s; ", f.Name, f.Value.String()) })
			name := DatabaseName(cmd)
			fmt.Println(name)
			fmt.Println()
			fmt.Println()
			if name != "splicedb" {
				t.Fatalf("database name shoud have been: 'splicedb'. but instead was: '%s'", name)
			}
		},
	}
	tcmd.Flags().StringP("database-name", "d", "", "")
	tcmd.Flags().String("database", "", "")
	tcmd.Flags().String("workspace", "", "")
	tcmd.SetArgs(args)
	tcmd.Execute()
}
