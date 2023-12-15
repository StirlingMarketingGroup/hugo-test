import axios from "axios";

export default async function main() {
    console.log('Hello, world!');

    const contactForm = document.querySelector<HTMLFormElement>('#contact-form');
    contactForm?.addEventListener('submit', async e => {
        e.preventDefault();

        axios.post('/contacts', {
            "name": (contactForm.querySelector('#contact-first-name') as HTMLInputElement).value + " " + (contactForm.querySelector('#contact-last-name') as HTMLInputElement).value,
            "company": (contactForm.querySelector('#contact-organization') as HTMLInputElement).value,
            "email": (contactForm.querySelector('#contact-email') as HTMLInputElement).value,
            "phone": "+" + (contactForm.querySelector('#contact-phone') as HTMLInputElement).value,
            "content": "TEST submission please ignore",
            "referrer": window.location.href,
            "reCaptcha": "yes"
        }).then(() => {
            console.log('Success!');
            window.location.href = '/contact-success/'
        }).catch(err => console.error(err));
    });
}