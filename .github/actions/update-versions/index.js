const { writeFileSync, readFileSync } = require("fs");

function main(input) {
    const version = input('version')

    if (version) {
        const versions = JSON.parse(readFileSync('./versions.json'))
        versions.push({version: version})
        writeFileSync('./versions.json', JSON.stringify(versions, null, 2))
    }
}

function getInput(name) {
    return process.env[`INPUT_${name.replace(/ /g, '_').toUpperCase()}`] || '';
}

main(getInput);
