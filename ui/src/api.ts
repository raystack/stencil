
export const getNamespaces = () => fetch("/v1beta1/namespaces").then(res => res.json());
export const getSchemas = (namespace : string) => fetch(`/v1beta1/namespaces/${namespace}/schemas`).then(res => res.json());
export const getVersions = (namespace: string, schema : string) => fetch(`/v1beta1/namespaces/${namespace}/schemas/${schema}/versions`).then(res => res.json());
export const getLatestSchema = (namespace: string, schema : string) => fetch(`/v1beta1/namespaces/${namespace}/schemas/${schema}`);
export const getVersionedSchema = (namespace: string, schema : string, version : number) => fetch(`/v1beta1/namespaces/${namespace}/schemas/${schema}/versions/${version}`)
