import type {UserConfig} from '@commitlint/types';

const Configuration: UserConfig = {
    /*
   * Resolve and load @commitlint/config-conventional from node_modules.
   * Referenced packages must be installed
   */
    extends: ['@commitlint/config-conventional'],
    parserPreset: {
        parserOpts: {
            headerPattern: /^\[(\w*)]\s(\w*)(\(.+\))?:\s(.*)$/,
            headerCorrespondence: ['scope', 'type', 'ticket', 'subject'],
        }
    },
    /*
     * Any rules defined here will override rules from @commitlint/config-conventional
     */
    rules: {
        'scope': async () => {
            // TODO: Setup auto generate scopes based on NX Workspace
            // const graph = await createProjectGraphAsync({
            //     exitOnError: true
            // })
            // const projectMap = await createProjectFileMapUsingProjectGraph(graph)
            return [2, 'always', ['workspace', 'api', 'web']]
        },
        'subject-case': [2, 'always', 'sentence-case'],
    },
    plugins: [
        {
            rules: {
                'scope': ({ scope }, _when, value: string[] = []) => {
                    return [
                        value.includes(scope || ''),
                        `[scope] should contains invalid Nx project name in this workspace (${value})`
                    ]
                }
            }
        }
    ],
};

module.exports = Configuration;