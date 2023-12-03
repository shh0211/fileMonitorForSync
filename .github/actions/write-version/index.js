const { writeFileSync, readFileSync } = require("fs");

function main(input) {
    const version = JSON.parse(readFileSync('./version.json')).version
    const changelog = readFileSync('./CHANGELOG.md')
    const gofile = `package version
var Version = "${version}"
var Changelog = ${JSON.stringify(changelog.toString())}
   `

    writeFileSync('./version/version.go', Buffer.from(gofile))
}

main();
