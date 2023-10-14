import type {UserConfig} from '@commitlint/types';
import {RuleConfigSeverity} from '@commitlint/types';

const Configuration: UserConfig = {
    parserOpts: {
        headerPattern: /^(\w*)\((\w*)\)-(\w*)\s(.*)$/,
        headerCorrespondence: ['type', 'scope', 'ticket', 'subject'],
    },
    /*
     * Resolve and load @commitlint/format from node_modules.
     * Referenced package must be installed
     */
    formatter: '@commitlint/format',
    /*
     * Any rules defined here will override rules from @commitlint/config-conventional
     */
    rules: {
        'type-enum': [RuleConfigSeverity.Error, 'always', ['foo']],
    },
    /*
     * Functions that return true if commitlint should ignore the given message.
     */
    ignores: [(commit) => commit === ''],
    /*
     * Whether commitlint uses the default ignore rules.
     */
    defaultIgnores: true,
};

module.exports = Configuration;