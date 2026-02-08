import { localize, setup } from '/static/lib/jinya-alpine-tools.js';

document.addEventListener('DOMContentLoaded', async () => {
  const MessagesDe = await (await fetch('/static/admin/langs/messages.de.json')).json();
  const MessagesEn = await (await fetch('/static/admin/langs/messages.en.json')).json();

  await setup({
    defaultArea: 'apps',
    defaultPage: '',
    baseScriptPath: '/static/admin/js/',
    routerBasePath: '/admin',
    storagePrefix: '/jinya/releases',
    openIdConfig: jinyaOpenIdConfig,
    languages: { de: MessagesDe, en: MessagesEn },
  });
});
