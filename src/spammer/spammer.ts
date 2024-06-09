import { GaxiosError, GaxiosResponse, request } from "gaxios";

import { logComplete, logDownload, logError } from "../cli/logger";
import { getConfig } from "../config";
import { NpmjsResponse } from "../models/npmjs-response.model";
import { Stats } from "../models/stats.model";
import { stripOrganisationFromPackageName } from "./utils";

export const queryNpms = async (): Promise<NpmjsResponse> => {
    const npmsResponse: GaxiosResponse<NpmjsResponse> = await request<NpmjsResponse>({
        baseUrl: "https://registry.npmjs.com",
        url: '/-/v1/search',
        params: {
          text: getConfig().packageName,
          size: 1,
        },
        method: "GET",
    })
    .then(((response) => {
        if (!response.data?.objects?.length) {
            throw Error(`Not found package`);
        }

        return response;
    }))
    .catch((response: GaxiosError<NpmjsResponse>) => {
        throw Error(`Failed to download.\n${response.message}`);
    });

    return npmsResponse.data;
};

export const downloadPackage = async (version: string, stats: Stats): Promise<unknown> => {
    const packageName: string = getConfig().packageName;
    const unscopedPackageName: string = stripOrganisationFromPackageName(packageName);

    return request<unknown>({
        baseUrl: "https://registry.yarnpkg.com",
        url: `/${packageName}/-/${unscopedPackageName}-${version}.tgz`,
        method: "GET",
        timeout: getConfig().downloadTimeout,
        responseType: "stream",
    })
        .then((_) => stats.successfulDownloads++)
        .catch((_) => stats.failedDownloads++);
};

const spamDownloads = async (version: string, stats: Stats): Promise<void> => {
    const requests: Promise<unknown>[] = [];

    for (let i = 0; i < getConfig().maxConcurrentDownloads; i++) {
        requests.push(downloadPackage(version, stats));
    }

    await Promise.all(requests);

    if (stats.successfulDownloads < getConfig().numDownloads) {
        await spamDownloads(version, stats);
    }
};

export const run = async (): Promise<void> => {
    try {
        const npmsResponse: NpmjsResponse = await queryNpms();
        const version: string = npmsResponse.objects[0].package.version;
        const startTime = Date.now();
        const stats: Stats = new Stats(startTime);
        const loggingInterval: NodeJS.Timeout = setInterval(() => logDownload(stats), 1000);
        await spamDownloads(version, stats);

        clearInterval(loggingInterval);
        logComplete();
    } catch (e) {
        logError(e);
    }
};
