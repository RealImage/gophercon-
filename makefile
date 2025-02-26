APP_NAME = qubecli
BIN_DIR = /usr/local/bin

.PHONY: all build install uninstall clean

# Default target (just runs build)
all: build

# Build the Go binary
build:
	go build -o $(APP_NAME)

# Install the binary system-wide
install: build
	sudo mv $(APP_NAME) $(BIN_DIR)/
	sudo chmod +x $(BIN_DIR)/$(APP_NAME)
	echo "Installation complete! Run '$(APP_NAME)' from anywhere."

# Remove the installed binary
uninstall:
	sudo rm -f $(BIN_DIR)/$(APP_NAME)
	echo "Uninstalled $(APP_NAME)."

# Clean up build artifacts
clean:
	rm -f $(APP_NAME)
	echo "Cleaned up build files."
