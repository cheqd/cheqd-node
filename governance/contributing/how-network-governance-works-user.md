# Minor Network changes

## **Minor Network changes**

These are changes that are insignificant, and do not affect the way the Network runs overall. Minor Network changes include, but are not limited to:

* Textual edits to the Governance Framework or written general documentation;
* Sensible additions to general documentation to improve clarity;
* Minor code changes;
* Finding, reporting and suggesting fixes to bugs;
* Translating the cheqd documentation into various languages.

## Decision tree

To help YOU understand how to make changes on the cheqd Network, the decision tree below visualises how changes should be carried out.

![Decision tree for cheqd Governance](../../.gitbook/assets/on-chain-vs-off-chain-decision-tree-1-.jpg)

## How do I make edits to general text or code in cheqd documentation?

You are able to make revisions and amendments to the wiki and source code without making a formal Proposal. 

cheqd is an Open Source project which means that anyone is free to propose a revision, addition or amendment. 

Changes can be made through two routes:

1. [GitBook](https://docs.cheqd.io/cheqd-node/)
2. [Github](https://github.com/cheqd)

### GitBook

GitBook is where cheqd's Wiki lives and where YOU can make suggested changes to the written documentation. 

To make a change, you need to use the link below to become an Open Source Contributor on cheqd's GitBook:

{% embed url="https://app.gitbook.com/invite/cheqd?invite=-MiVxCDUlLSB22RuQ6dl" %}

You should follow these instructions to make a change to any cheqd GitBook content.

1. **Make your own branch**

To create a new branch, select 'New' followed by 'New variant'.

![](../../.gitbook/assets/image%20%281%29.png)

2. Next you need to **name your new branch**.

You should use the prefix:

* ADR/
* gov/

followed by a unique title indicating your draft number such as

* ADR/{username}draft1
* gov/{username}draft1

On this branch you should make and save all desired changes to the content.

Once you have finished your changes, save the branch but **do not merge the branch to main**. 

Your edits will be reviewed by a cheqd admin or moderator and then merged, amended or rejected.

### Github

You may also use Github to make changes and amendments to both the source code and the written content in this documentation.

You should similarly make a branch of the cheqd Github main. 

You should then make a pull request with your proposed changes, edits, revisions and additions. 

**Please use clean, concise titles for your pull requests.** Assume that the person reading the pull request title is not a programmer, but instead a cheqd Network user or lay-person, and **try to describe your change or fix from in plain English**. We use commit squashing, so the final commit in the main branch will carry the title of the pull request, and commits from the main branch are fed into the changelog. The changelog is separated into [keepachangelog.com categories](https://keepachangelog.com/en/1.0.0/), and while that spec does not prescribe how the entries ought to be named, for easier sorting, start your pull request titles using one of the verbs "Add", "Change", "Deprecate", "Remove", or "Fix" \(present tense\).

Example:

| Not ideal | Better |
| :--- | :--- |
| Fixed NoMethodError in RemovalWorker | Fix nil error when removing statuses caused by race condition |

It is not always possible to phrase every change in such a manner, but it is desired.

**The smaller the set of changes in the pull request is, the quicker it can be reviewed and merged.** Splitting tasks into multiple smaller pull requests is often preferable.

**Pull requests that do not pass automated checks may not be reviewed**. 

## Bug reports

During the course of cheqd's lifecycle, it is likely that bugs will arise. Reporting these bugs effectively and resolving them promptly will be very important for the smooth running of the Network.

Bugs SHOULD be discussed or raised in the [Technical Help section of our Forum](https://github.com/cheqd/cheqd-node/discussions/categories/technical-help). 

From a bug raised here, bugs SHOULD be forwarded to [GitHub Issues](https://github.com/cheqd/cheqd-node/issues). This can be done by clicking the button "create issue from discussion".

Please use the search function to make sure that you are not submitting duplicates, and that a similar report or request has not already been resolved or rejected.

## Translations

You can submit translations via Github or GitBook branches \(as above\). We would greatly appreciate this work from our community to ensure that cheqd reaches a diverse, multijurisdictional and multilingual spread of society!  
 

## 

##  

