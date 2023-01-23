# Test Website

Welcome to our test website!

This site is assembled very similarly to how our real websites are built. One of our sites, for example, is <https://my.allaboutchallengecoins.com/login/>.

This website is put together with a couple of tools to help make writing and maintaining it easier, including:

- [Hugo](https://gohugo.io/) - The framework we use to help us organize our page content and HTML, and enables us to live-edit our websites while we write them.
- [Bootstrap 5.1](https://getbootstrap.com/docs/5.1/getting-started/introduction/) - SCSS framework for styling everything about our websites, making responsive website building easy, and making the look across our website consistent.
- [TypeScript](https://www.typescriptlang.org/) - The language used to write all client side functionality. No react or angular here, just plain old TypeScript. This example repo adds shortcuts to VS Code to help us compile our TypeScript.

These three tools all have very thorough and easy to use documentation, and it's the core of how we build all of our sites.

This project contains only client-side (front end) code. It is the goal of this client side code to make it as easy as possible for our customers to interact with our system. Customers come to our websites to request quotes for their custom products, add their quotes to their cart, order their carts, see order history, and more. We'd like to expand this functionality to include group pay, custom product builders, and all sorts of fun stuff.

You don't have to worry about setting up the project like this one initially, but you will be responsible for adding features and updating the files in this project.

## The test

This front end code will interact with our own, internally written, REST API, which is responsible for taking the customer inputs you give us and processing that data.

In this project, there is a contact page started that needs some blanks filled in.

- [content/contact.html](content/contact.html) has a form element but no fields.
    - Use <https://getbootstrap.com/docs/5.1/forms/layout/> to add the relevant contact form fields to the form. This should include a first name, last name, organization name, email address, and phone number.
- [assets/ts/pages/contact.ts](assets/ts/pages/contact.ts) has a call to our API using [axios](https://www.npmjs.com/package/axios), but it isn't finished yet. The call should POST all the fields you added in the previous step to our API.
    - Our API has its own documentation that you can find [here](https://api2.allaboutchallengecoins.com/docs/#/Contact%20Form%20Submissions/post_contacts).

Once the call succeeds, the client should redirect the user to a "Contact Success" page, which doesn't exist yet. It should be created as [content/contact-success.html](content/contact-success.html) and should say something like:

> Thank you for your request! We will reach out within 24 hours!

The API that this is connected to is *live*! So don't send anything too stupid. Our sales people are used to ignoring test messages.

## VS Code

We use VS Code shortcuts to help build and test our websites. The two things you need to know for this test are:

1. Press Ctrl+Shift+B to bring up the build menu, and select and run the entry titled "Compile TS" to compile the TypeScript. You will see the compiling status of your code in the VS Code Terminal, and if it succeeds, the Hugo Test site will refresh with your updated, compiled code.
2. Press Ctrl+Shift+B to bring up the build menu, and select and run the entry titled "Hugo Server: Test Site". Once this is running, you should be able to navigate in your browser to <http://test.localhost:10506/contact/> to see your cool new contact form.

You can view your current running tasks in VS Code by clicking this icon at the bottom:

![Show Running Tasks](images/showbuildtasks.png)

If you are familiar with GitHub, creating a fork of this repository with your changes would be fantastic.

---

Bonus! The website's styles are written in SCSS, which resides in the file [assets/scss/app.scss](assets/scss/app.scss). Edit this file to make this website whatever color you want!

---

If you have any questions *at all*, just ask them!

## Installing Prerequisites (Windows)

If I didn't set up the environment for you, there are a few things we need to do before the above VS Code short cuts will work. For all the installers below, simply accept all the defaults and everything should be gravy.

1. Install Golang - <https://go.dev/dl/go1.19.5.windows-amd64.msi>
2. Install hugo - <https://github.com/gohugoio/hugo/releases/download/v0.110.0/hugo_extended_0.110.0_windows-amd64.zip>
    1. Open up explorer.exe.
    2. Navigate to %userprofile% and create directories "go", and then inside "go" create a "bin" directory.
    3. Copy "hugo.exe" from the hugo zip to this new "bin" folder (should be %userprofile%\go\bin)
3. [Download & Install Git](https://github.com/git-for-windows/git/releases/download/v2.39.1.windows.1/Git-2.39.1-64-bit.exe). This one asks a lot of questions, but just all the defaults are fine.
4. [Download & Install "nvm-windows"](https://github.com/coreybutler/nvm-windows/releases/download/1.1.10/nvm-setup.exe)
    1. Open up Powershell and run `nvm install latest`
    2. Next, run `nvm use latest`
    3. To ensure this worked, run `node --version`, and you should see the version of nodejs you just installed
    4. Install all of the following nodeJS prerequisites for this project by running `npm i -g google-closure-compiler tslib ts yarn`
5. [Install VS Code](https://code.visualstudio.com/download) (if not installed already, of course)
6. [Install GitHub Desktop](https://desktop.github.com/).
    1. Create an account or log in, and choose "Clone a repository from the internet", then click on the "URL" tab on the right.
    2. First, download the repo "https://github.com/BrianLeishman/tscc". It might ask you to trust the repo, accept if you do :)
        1. Right click on top-left where it says "Current repository" and "tscc" and click "Copy repo path"
        2. In Powershell, type "dir" followed by a space, and then paste from the clipboard, and hit enter. Your current folder should look like "C:\Users\brian\Documents\GitHub\tscc"
        3. Next, run all of the following commands:
            - `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Unrestricted` (then press A, ENTER)
            - `yarn install`
            - `yarn lerna clean --yes`
            - `yarn lerna bootstrap` (Allow the network connection if asked)
            - `yarn tsc`
            - `cd packages/tscc`
            - `npm remove -g .`
            - `npm i -g .`
    3. Next, click File -> Clone Repository and clone "https://github.com/StirlingMarketingGroup/hugo-test". Then press Ctrl+Shift+A to open this repo in VS Code, and say yes to trusting the author
        1. From inside VS Code, press Ctrl+Shift+B and run "Node Install". Look in the running scripts to see the status.
        2. Next, try running Compile TS, and hopefully it works! If it does, you're set up :)
        3. Refer to the above instructions to run the Hugo site locally and view it in the browser.
