# Solar Garlic Mailing List

Mailing list backend for the Solar Garlic Band Website.

## Design

### Registration Flow Happy Path

1. User submits registration form in FE
    - request to /email/subscribe endpoint
1. Backend generates a registration code with a 24hr expiration, writes to database
    - token is a string with 64-128 randomly generated characters
1. User is redirected to a verification email sent page
1. Backend sends an email to the provided address with a registration link containing the registration code
1. User clicks the registration link in the email. Backend writes the email to the mailing list and redirects the user to a success page.
    - links to /email/verify endpoint
1. User receives an email telling them they were subscribed


### Send Email Flow Happy Path

1. Admin drafts an email message and submits to a mailing list
1. Every user on the mailing list receives the email message
1. Unsubscribe link in each email must contain the token for the user being emailed


### Unsubscribe Flow

1. User clicks unsubscribe link at the bottom of an email sent to the mailing list. Link needs to contain a way to uniquely identify the user. Reuse the random token that was generated on account creation.
1. User is redirected to an unsubscribed page on the frontend
1. User receives an email telling them they have been unsubscribed
