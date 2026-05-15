// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightLlmsTxt from 'starlight-llms-txt';
import starlightLinksValidator from 'starlight-links-validator';
import mermaid from 'astro-mermaid';

// https://astro.build/config
export default defineConfig({
site: 'https://githubnext.github.io',
base: '/apm/',
integrations: [
mermaid(),
starlight({
title: 'APM Go Migration Progress',
description: 'Current status, benchmark signals, and next work for the APM Python-to-Go migration.',
favicon: '/favicon.svg',
social: [
{ icon: 'github', label: 'GitHub', href: 'https://github.com/githubnext/apm' },
],
tableOfContents: {
minHeadingLevel: 2,
maxHeadingLevel: 4,
},
pagination: false,
customCss: ['./src/styles/custom.css'],
expressiveCode: {
frames: {
showCopyToClipboardButton: true,
},
},
plugins: [
starlightLinksValidator({
errorOnRelativeLinks: false,
errorOnLocalLinks: true,
}),
starlightLlmsTxt({
description: 'Current status, benchmark signals, and next work for the APM Python-to-Go migration.',
}),
],
sidebar: [
{
label: 'Progress',
items: [
{ label: 'Autoloop Go Migration', slug: 'index' },
],
},
],
}),
],
});
