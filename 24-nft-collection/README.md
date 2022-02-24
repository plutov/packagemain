## packagemain #24: Generate an NFT Collection in Go

Nowadays almost any conversation will inevitably end up with discussing NFTs or non-fungible tokens. However, I'm not going to talk about it much today, as I am not an expert. I will rather focus on a technical part of generating a collection of unique images by blending the layers in Go. That should be fun. And probably later if I gain more knowledge in this space we can "mint" them as well.

### How images will be generated?

To generate the images for a collection, we'll be using mutiple types of layers and blend them programmatically. For this example I decided to use the following 3 layers:
- 10 background images
- 10 gophers
- 10 quotes that gophers say

This will allow us to create 1000 unique images for our collection. Each emage will contain only 1 background, 1 gopher and 1 quote.

I am not an artist, so I had to find some assets to do this video. I want to expect the authors' licenses, so I carefully selected few assets that we can legally use here.

- Backgrounds: https://unsplash.com/@cocoloris (license: https://unsplash.com/license)
- Gophers: https://github.com/MariaLetta/free-gophers-pack (license: https://github.com/MariaLetta/free-gophers-pack/blob/master/LICENSE)
- Quotes: I generated myself 10 simple quotes

I already uploaded all images that we'll use into the separate folders.

### What the program will do

Before writing the program, let's define what it'll do. We will build an executable program in which we'll define where to find layers' images. The prohram then will blend all unique combinations programmatically, and store them into a separate folder.

This program will rely a lot on [image](https://pkg.go.dev/image) package.

### Define types and functions

Now, as everything is ready, let's define some basic types and functions.
