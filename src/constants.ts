import figlet from "figlet";
import standard from "figlet/importable-fonts/Standard.js";
import colors from "picocolors";

figlet.parseFont("Standard", standard);

const { magenta, cyan, yellow, green } = colors;

export const LOGO = magenta(
  figlet.textSync("SpecLedger", { font: "Standard", horizontalLayout: "full" })
);

export const helpMessage = `
${LOGO}

${cyan("Usage:")} specledger [command]

${yellow("Commands:")}
  ${green("init")}     Initialize SpecLedger in current project
  ${green("help")}     Show this help message

${yellow("Examples:")}
  npx specledger          Interactive CLI
  npx specledger init     Initialize SpecLedger files
`;
