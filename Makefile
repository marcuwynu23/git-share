BINARY   := git-share
MODULE   := github.com/markwayne/git-share
GO       := go

UNAME_S  := $(shell uname -s 2>/dev/null || echo "windows")

ifeq ($(UNAME_S),Linux)
	OUTPUT    := $(BINARY)
	RM        := rm -f
	CP        := cp -f
	RUN       := ./
	INSTALL   := /usr/local/bin
endif
ifeq ($(UNAME_S),Darwin)
	OUTPUT    := $(BINARY)
	RM        := rm -f
	CP        := cp -f
	RUN       := ./
	INSTALL   := /usr/local/bin
endif
ifneq (,$(findstring CYGWIN,$(UNAME_S)))
	OUTPUT    := $(BINARY).exe
	RM        := rm -f
	CP        := cp -f
	RUN       := ./
	INSTALL   := /usr/local/bin
endif
ifneq (,$(findstring MINGW,$(UNAME_S)))
	OUTPUT    := $(BINARY).exe
	RM        := rm -f
	CP        := cp -f
	RUN       := ./
	INSTALL   := /usr/local/bin
endif
ifneq (,$(findstring MSYS,$(UNAME_S)))
	OUTPUT    := $(BINARY).exe
	RM        := rm -f
	CP        := cp -f
	RUN       := ./
	INSTALL   := /usr/local/bin
endif
ifeq ($(UNAME_S),windows32)
	OUTPUT    := $(BINARY).exe
	RM        := cmd /c del /f
	CP        := cmd /c copy /y
	RUN       :=
	INSTALL   := C:\Bin\tools
endif
ifeq ($(UNAME_S),windows)
	OUTPUT    := $(BINARY).exe
	RM        := cmd /c del /f
	CP        := cmd /c copy /y
	RUN       :=
	INSTALL   := C:\Bin\tools
endif

.PHONY: all build clean test vet lint run install help

all: build

build:
	$(GO) build -o $(OUTPUT) ./cmd/$(BINARY)

clean:
	$(GO) clean
	-$(RM) $(OUTPUT) 2>nul
	-$(RM) $(BINARY) 2>nul
	-$(RM) $(BINARY).exe 2>nul

test:
	$(GO) test -v -race -count=1 ./...

vet:
	$(GO) vet ./...

lint:
	golangci-lint run ./...

run: build
	$(RUN)$(OUTPUT)

install: build
	$(CP) $(OUTPUT) $(INSTALL)\$(OUTPUT)

help:
	@echo "Targets:"
	@echo "  build   - Build the $(BINARY) binary"
	@echo "  clean   - Remove build artifacts"
	@echo "  test    - Run tests with race detector"
	@echo "  vet     - Run go vet"
	@echo "  lint    - Run golangci-lint"
	@echo "  run     - Build and run"
	@echo "  install - Build and copy to $(INSTALL)"
