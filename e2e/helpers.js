// @ts-check
const { spawn } = require("child_process");
const fs = require("fs");
const path = require("path");
const os = require("os");

/**
 * Creates a temporary directory with markdown files for testing
 * @param {Object.<string, string>} files - Map of filename to content
 * @returns {string} Path to the temp directory
 */
function createTestSite(files) {
  const testDir = fs.mkdtempSync(path.join(os.tmpdir(), "volcano-e2e-"));

  for (const [filename, content] of Object.entries(files)) {
    const filePath = path.join(testDir, filename);
    const dir = path.dirname(filePath);
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }
    fs.writeFileSync(filePath, content);
  }

  return testDir;
}

/**
 * Starts a volcano server
 * @param {string} testDir - Directory to serve
 * @param {number} port - Port number
 * @param {string[]} [extraArgs] - Additional CLI arguments
 * @param {Object} [options] - Options
 * @param {boolean} [options.noDefaults] - Don't add default --instant-nav and --search flags
 * @returns {Promise<import('child_process').ChildProcess>}
 */
async function startServer(testDir, port, extraArgs = [], options = {}) {
  const defaultFlags = options.noDefaults ? [] : ["--instant-nav", "--search"];
  const args = [
    "serve",
    "-p",
    String(port),
    ...defaultFlags,
    ...extraArgs,
    testDir,
  ];

  const serverProcess = spawn("./volcano", args, {
    cwd: path.join(__dirname, ".."),
    stdio: "pipe",
  });

  // Wait for server to be ready
  await new Promise((resolve, reject) => {
    const timeout = setTimeout(
      () => reject(new Error("Server startup timeout")),
      10000,
    );

    serverProcess.stdout.on("data", (data) => {
      if (data.toString().includes("http://localhost")) {
        clearTimeout(timeout);
        resolve(undefined);
      }
    });

    serverProcess.stderr.on("data", (data) => {
      // Log errors but don't fail - some warnings are expected
      if (process.env.DEBUG) {
        console.error("Server stderr:", data.toString());
      }
    });

    serverProcess.on("error", (err) => {
      clearTimeout(timeout);
      reject(err);
    });
  });

  return serverProcess;
}

/**
 * Stops a volcano server and cleans up
 * @param {import('child_process').ChildProcess} serverProcess
 * @param {string} testDir
 */
function stopServer(serverProcess, testDir) {
  if (serverProcess) {
    serverProcess.kill("SIGTERM");
  }
  if (testDir && fs.existsSync(testDir)) {
    fs.rmSync(testDir, { recursive: true, force: true });
  }
}

/**
 * Sample content for generating longer pages
 */
const LOREM = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.`;

module.exports = {
  createTestSite,
  startServer,
  stopServer,
  LOREM,
};
