import axios from "axios";

export default async function main() {
    console.log('Hello, world!');

    const contactForm = document.querySelector<HTMLFormElement>('#contact-form');
    contactForm?.addEventListener('submit', async e => {
        e.preventDefault();

        axios.post('/contact', {
            // params and what
        }).then(resp => {
            console.log('Success!');
            // redirect to success page
        }).catch(err => alert(err));
    });
}