#!/usr/bin/env node
import * as prompts from "@clack/prompts";
import mri from "mri";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { execSync } from "node:child_process";
import colors from "picocolors";
import { LOGO, helpMessage } from "./constants";
import { copyDir } from "./utils";

const { green, cyan, yellow } = colors;

async function main() {
  const argv = mri<{
    help?: boolean;
  }>(process.argv.slice(2), {
    alias: { h: "help" },
    boolean: ["help"],
  });

  const command = argv._[0];

  // Show help
  if (argv.help || command === "help") {
    console.log(helpMessage);
    return;
  }

  // Direct init command
  if (command === "init") {
    await handleInit();
    return;
  }

  // Interactive CLI
  await showInteractiveCLI();
}

async function showInteractiveCLI() {
  console.log();
  console.log(LOGO);
  console.log();

  prompts.intro(cyan("Welcome to SpecLedger CLI"));

  const action = await prompts.select({
    message: "What would you like to do?",
    options: [
      {
        value: "init",
        label: green("Initialize SpecLedger"),
        hint: "Copy SpecLedger files to current project",
      },
      {
        value: "help",
        label: yellow("Show help"),
        hint: "Display available commands",
      },
    ],
  });

  if (prompts.isCancel(action)) {
    prompts.cancel("Operation cancelled");
    return;
  }

  if (action === "init") {
    await handleInit();
  } else if (action === "help") {
    console.log(helpMessage);
  }
}

function isBeadsInstalled(): boolean {
  try {
    execSync("which beads", { stdio: "ignore" });
    return true;
  } catch {
    return false;
  }
}

async function installBeads(): Promise<boolean> {
  const spinner = prompts.spinner();
  spinner.start("Installing beads...");

  try {
    execSync(
      "curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash",
      { stdio: "inherit" }
    );
    spinner.stop(green("Beads installed successfully!"));
    return true;
  } catch (error) {
    spinner.stop("Failed to install beads");
    prompts.log.error(
      `Error: ${error instanceof Error ? error.message : "Unknown error"}`
    );
    return false;
  }
}

async function handleInit() {
  const cwd = process.cwd();
  const initDir = path.resolve(
    fileURLToPath(import.meta.url),
    "..",
    "init"
  );

  prompts.intro(cyan("Initializing SpecLedger..."));

  // Check if beads is installed
  if (!isBeadsInstalled()) {
    prompts.log.warn(yellow("Beads is not installed."));

    const shouldInstall = await prompts.confirm({
      message: "Would you like to install beads now?",
    });

    if (prompts.isCancel(shouldInstall)) {
      prompts.cancel("Operation cancelled");
      return;
    }

    if (shouldInstall) {
      const success = await installBeads();
      if (!success) {
        prompts.log.error("Please install beads manually and try again.");
        process.exit(1);
      }
    } else {
      prompts.log.warn("Beads is required for SpecLedger. Please install it manually:");
      prompts.log.info(cyan("curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash"));
      process.exit(1);
    }
  }

  const spinner = prompts.spinner();
  spinner.start("Copying SpecLedger files...");

  try {
    copyDir(initDir, cwd);
    spinner.stop("Files copied successfully!");

    prompts.outro(
      green("SpecLedger initialized!") +
        "\n\n" +
        "  Next steps:\n" +
        cyan("  1.") + " Review the copied files\n" +
        cyan("  2.") + " Run " + yellow("chmod +x setup.sh && ./setup.sh") + " to complete setup\n"
    );
  } catch (error) {
    spinner.stop("Failed to copy files");
    prompts.log.error(
      `Error: ${error instanceof Error ? error.message : "Unknown error"}`
    );
    process.exit(1);
  }
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
