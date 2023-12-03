const { readFileSync } = require("fs");

async function main(output) {
    output('version', JSON.parse(readFileSync('./version.json')).version)
    output('changelog', readFileSync('./CHANGELOG.md'))
    output('current_date', new Date().toLocaleDateString().split('/').join('-'))
}

function setOutput(name, value) {
    process.stdout.write(Buffer.from(`::set-output name=${name}::${value}\n`))
}

main(setOutput);
