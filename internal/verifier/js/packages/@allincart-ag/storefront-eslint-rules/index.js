import MigratePluginManager from './plugin-manager.js';
import DomAccessHelper from "./dom-access-helper.js";
import HttpClient from "./http-client.js";
import QueryString from "./query-string.js";

export default {
    plugins: {
        "allincart-storefront": {
            rules: {
                "migrate-plugin-manager": MigratePluginManager,
                "no-dom-access-helper": DomAccessHelper,
                "no-http-client": HttpClient,
                'no-query-string': QueryString,
            },
        }
    },
    rules: {
        'allincart-storefront/migrate-plugin-manager': 'error',
        'allincart-storefront/no-dom-access-helper': 'error',
        'allincart-storefront/no-http-client': 'error',
        'allincart-storefront/no-query-string': 'error',
    }
}