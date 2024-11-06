#!/usr/bin/env node

import { setConfig } from "./config";
import { Config } from "./models/config.model";
import { run } from "./spammer/spammer";

const config: Config = {
    packageName: process.env.NPM_PACKAGE_NAME as string,
    numDownloads: process.env.NPM_NUM_DOWNLOADS as unknown as number || 1000,
    maxConcurrentDownloads: process.env.NPM_MAX_CONCURRENT_DOWNLOAD as unknown as number || 300,
    downloadTimeout: process.env.NPM_DOWNLOAD_TIMEOUT as unknown as number || 3000
}

setConfig(config);
run();
