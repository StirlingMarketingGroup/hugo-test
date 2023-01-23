import Contact from './pages/contact';

async function main() {
    switch (document.body.id) {
        case 'contact':
            Contact();
            break;
    }
}

main();