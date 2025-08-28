# Easy Commits

An AI-powered git commit message generator that creates detailed, conventional commit messages based on your git diff and optional user context.

## Features

- **Multiple AI Providers**: Supports OpenAI, Anthropic Claude, and local Ollama models
- **Conventional Commits**: Generates messages following conventional commit format
- **User Context**: Add additional context to help the AI understand your changes
- **Interactive**: Shows generated message and asks for confirmation before committing
- **Smart Diff Detection**: Automatically detects staged or unstaged changes

## Installation

1. Clone this repository
2. Build the binary:
   ```bash
   go build -o easy-commits
   ```
3. (Optional) Move to your PATH:
   ```bash
   sudo mv easy-commits /usr/local/bin/
   ```

## Setup

First, configure your AI provider:

```bash
easy-commits config
```

You'll be prompted to select a provider and enter your credentials:

### OpenAI
- Provider: `openai`
- API Key: Your OpenAI API key
- Model: `gpt-3.5-turbo` (default)

### Anthropic Claude
- Provider: `anthropic`
- API Key: Your Anthropic API key
- Model: `claude-3-haiku-20240307` (default)

### Ollama (Local)
- Provider: `ollama`
- Base URL: `http://localhost:11434` (default)
- Model: Your local model name (e.g., `llama2`, `codellama`)

Configuration is saved to `~/.easy-commits-config.json`

## Usage

### Basic Usage

Stage your changes and generate a commit:
```bash
git add .
easy-commits commit
```

Or let easy-commits stage all changes automatically:
```bash
easy-commits commit
```

### With Additional Context

Provide extra context to help the AI understand your changes:
```bash
easy-commits commit --context "Fixed the authentication bug that was causing login failures"
```

### Example Workflow

1. Make your changes
2. Run `easy-commits commit`
3. Review the generated commit message
4. Confirm or cancel the commit

## Generated Commit Message Format

The tool generates commit messages following the conventional commit format:

```
type(scope): description

Optional longer description explaining the changes
```

Types include:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Maintenance tasks

## Commands

- `easy-commits config` - Configure AI provider and credentials
- `easy-commits commit` - Generate and create a commit
- `easy-commits commit --context "description"` - Add user context
- `easy-commits help` - Show help information

## Requirements

- Go 1.24.6 or later
- Git repository
- API key for chosen AI provider (or local Ollama setup)

## License

MIT License