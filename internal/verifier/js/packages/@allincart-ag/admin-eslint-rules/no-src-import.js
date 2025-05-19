export default {
    create(context) {
        return {
            ImportDeclaration(node) {
                const invalidNodeSources = [];
                invalidNodeSources.push(node.source.value.startsWith('@administration/'));

                if (invalidNodeSources.includes(true)) {
                    context.report({
                        loc: node.source.loc.start,
                        message: `\
You can't use imports directly from the Allincart Core via "${node.source.value}". \
Use the global Allincart object directly instead (https://developer.allincart.net/docs/guides/plugins/plugins/administration/the-allincart-object)`,
                    });
                }
            },
        };
    },
};
