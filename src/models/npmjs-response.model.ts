interface NpmjsPackage {
  name: string;
  scope: string;
  version: string;
  description: string;
  keywords: string[];
  date: string;
  links: Record<string, string>;
  publisher: Record<string, string>;
  maintenaners: Record<string, string>;
}

interface NpmjsObject {
    package: NpmjsPackage;
}

export interface NpmjsResponse {
    objects: NpmjsObject[];
}
