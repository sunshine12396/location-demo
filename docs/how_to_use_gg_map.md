> **📌 NOTE: This project uses OpenStreetMap (OSM) by default.**
> You do **not** need a Google API key or billing account for the application to work out of the box. The service is pre-configured with a free Nominatim client. Only follow this guide if you wish to switch back to the `GooglePlacesClient`.

Phase 1: Setup Your API Key
Before you can write code, you need to register with the Google Cloud Console.

Create a Project: Go to the Google Cloud Console. Click the project dropdown at the top and select New Project. Give it a name like "My Maps Project."

Enable Billing: Even though there is a generous free tier (usually $200 in monthly credits), Google requires a linked billing account (credit card or bank account) to prevent bot abuse.

Enable APIs:

In the sidebar, go to APIs & Services > Library.

Search for the specific API you need (e.g., Geocoding API, Maps JavaScript API, or Places API).

Click on it and hit Enable.

Generate the Key:

Go to APIs & Services > Credentials.

Click + Create Credentials at the top and select API key.

Crucial Step: Click Edit Settings on your new key. Under "API restrictions," limit the key to only the APIs you enabled. Under "Application restrictions," restrict it to your website's URL or your IP address so others can't steal your credits.

Phase 2: Testing with cURL
The Geocoding API is the easiest to test via cURL. It converts an address (like "1600 Amphitheatre Parkway") into geographic coordinates.

The Code Template
Replace YOUR_API_KEY with the key you just generated.

Bash
curl -X GET "https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key=YOUR_API_KEY"
Breaking Down the Command
curl -X GET: Tells your terminal to "get" data from a URL.

json: This specifies that you want the response formatted in JSON (the standard for web apps).

address=: The location you want to look up. Note that spaces are replaced with +.

key=: Your secure credential that authorizes the request.