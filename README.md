# spiderwoman
"Vertical" crawler, which main target is to count links (resolved, e.g. from bit.ly) to external domains from all pages of given resources

For example we have a website domain.com with index page and two other pages. On all the pages of domain.com there is a link http://goo.gl/blah which resolves to example.com. So the spiderwoman after full crawl of domain.com must get such result "example.com:3", that means 3 pages of domain.com links to example.com (and shortlink is not a problem, spiderwoman have to resolve it).
