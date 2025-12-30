.PHONY: test lint tools coverage

# テストを実行し、カバレッジを測定
test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# golangci-lintでコード検査を実行
lint:
	golangci-lint run

# 依存ツールをインストール
tools:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing octocov..."
	@go install github.com/k1LoW/octocov@latest

# octocovを使用してカバレッジレポートを生成・確認（80%以上をチェック）
coverage: test
	octocov

