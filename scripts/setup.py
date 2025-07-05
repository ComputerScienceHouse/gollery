import os


def main():
    # Change working directory to project root (parent of scripts)
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.abspath(os.path.join(script_dir, os.pardir))
    os.chdir(project_root)

    # Create directories
    required_dirs: list[str] = [
        ".container_volumes",
        ".container_volumes/gollery",
        ".container_volumes/gollery/logs",
        ".container_volumes/postgres",
        ".container_volumes/postgres/logs",
        ".container_volumes/postgres/data",
    ]

    for directory in required_dirs:
        os.makedirs(directory, exist_ok=True)

    # Copy example.env to .env if it doesn't exist
    env_file: str = ".env"
    example_env_file: str = "example.env"

    if not os.path.exists(env_file):
        if os.path.exists(example_env_file):
            with open(example_env_file, "r") as src:
                with open(env_file, "w") as dst:
                    dst.write(src.read())
        else:
            print(f"Warning: {example_env_file} not found. Skipping .env creation.")


if __name__ == "__main__":
    main()
