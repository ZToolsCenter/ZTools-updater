# 二进制文件名称
BINARY_NAME=ztools-updater

# 构建目录
BUILD_DIR=dist

.PHONY: all clean build-mac build-mac-amd64 build-mac-universal build-win

all: clean build-mac-universal build-win

build-mac:
	@echo "Building for macOS (ARM64)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/mac-arm64/$(BINARY_NAME) .

build-mac-amd64:
	@echo "Building for macOS (AMD64)..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/mac-amd64/$(BINARY_NAME) .

build-mac-universal: build-mac build-mac-amd64
	@echo "Creating macOS Universal Binary..."
	@mkdir -p $(BUILD_DIR)/mac-universal
	lipo -create -output $(BUILD_DIR)/mac-universal/$(BINARY_NAME) $(BUILD_DIR)/mac-arm64/$(BINARY_NAME) $(BUILD_DIR)/mac-amd64/$(BINARY_NAME)

build-win:
	@echo "Building for Windows (AMD64)..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/win-amd64/$(BINARY_NAME).exe .

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
