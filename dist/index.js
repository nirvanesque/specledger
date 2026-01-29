#!/usr/bin/env node
import*as e from"@clack/prompts";import t from"mri";import n from"node:path";import{fileURLToPath as r}from"node:url";import i from"picocolors";import a from"figlet";import o from"figlet/importable-fonts/Standard.js";import s from"node:fs";a.parseFont(`Standard`,o);const{magenta:c,cyan:l,yellow:u,green:d}=i,f=c(a.textSync(`SpecLedger`,{font:`Standard`,horizontalLayout:`full`})),p=`
${f}

${l(`Usage:`)} specledger [command]

${u(`Commands:`)}
  ${d(`init`)}     Initialize SpecLedger in current project
  ${d(`help`)}     Show this help message

${u(`Examples:`)}
  npx specledger          Interactive CLI
  npx specledger init     Initialize SpecLedger files
`;function m(e,t){let n=s.statSync(e);n.isDirectory()?h(e,t):s.copyFileSync(e,t)}function h(e,t){s.mkdirSync(t,{recursive:!0});for(let r of s.readdirSync(e)){let i=n.resolve(e,r),a=n.resolve(t,r);m(i,a)}}const{green:g,cyan:_,yellow:v}=i;async function y(){let e=t(process.argv.slice(2),{alias:{h:`help`},boolean:[`help`]}),n=e._[0];if(e.help||n===`help`){console.log(p);return}if(n===`init`){await x();return}await b()}async function b(){console.log(),console.log(f),console.log(),e.intro(_(`Welcome to SpecLedger CLI`));let t=await e.select({message:`What would you like to do?`,options:[{value:`init`,label:g(`Initialize SpecLedger`),hint:`Copy SpecLedger files to current project`},{value:`help`,label:v(`Show help`),hint:`Display available commands`}]});if(e.isCancel(t)){e.cancel(`Operation cancelled`);return}t===`init`?await x():t===`help`&&console.log(p)}async function x(){let t=process.cwd(),i=n.resolve(r(import.meta.url),`..`,`init`);e.intro(_(`Initializing SpecLedger...`));let a=e.spinner();a.start(`Copying SpecLedger files...`);try{h(i,t),a.stop(`Files copied successfully!`),e.outro(g(`SpecLedger initialized!`)+`

  Next steps:
`+_(`  1.`)+` Review the copied files
`+_(`  2.`)+` Run `+v(`chmod +x setup.sh && ./setup.sh`)+` to complete setup
`)}catch(t){a.stop(`Failed to copy files`),e.log.error(`Error: ${t instanceof Error?t.message:`Unknown error`}`),process.exit(1)}}y().catch(e=>{console.error(e),process.exit(1)});