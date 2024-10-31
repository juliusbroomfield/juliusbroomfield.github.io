+++
title = 'Projects'
date = 2024-04-14T15:03:42-04:00
menu = 'main'
weight = 50
draft = false
+++

Here are a few projects I've worked or am working on (with a lot of progress!), not exhaustive and is always ongoing.

---

### 🎶 Explainable Music Recommendation System
This is a project I’ve been working on on-and-off for quite a while. While using Spotify’s Discover Weekly, I noticed a couple of patterns:

1. I had already heard many recommended songs I had heard them before on platforms like TikTok or other social media.
2. Spotify occasionally suggested songs that sampled or even interpolated tracks I had already liked - I like listening to hip-hop and jazz, so this was often noticable.

These patterns aren’t necessarily issues (they’re probably net goods — one is a goal of social recommenders, and the other fits content-based recommendations), but I would appreciate having some explainations behind the recommendations. This encouraged me to think more about using LLMs in recommendation systems - not only to improve explainability but also to bring more nuance and context (both socially and in content) to recommendations

I’m currently going through papers on [LLMs in knowledge graphs](https://dl.acm.org/doi/pdf/10.1145/3616855.3635853)
 and their [applications in recommendation systems](https://arxiv.org/pdf/2308.10835), and coding as I go through them, with the goal of developing an explainable MRS.
- [Code](https://github.com/juliusbroomfield/xai-music-recommender)

---

### 💡 Ideate
Ideate is a web app. I built alongside a team of 4 for PennApps XXIV. Users can record videos of them and their team ideating, sketching high-level system design and simple whiteboard coding, and automatically turn it to functional code in a Github repository. Multimodal models weren't widely avaliable during this time, so this application looked to fill that gap. Once users uploaded their video, we picked up text from their video's audio, and captured frames from the video, using OCR to pick up any relevant text. This was fed into the OpenAI API, which would then output code in response to the video.

- [Code](https://github.com/juliusbroomfield/ideate)

---

### ♟️ Chess Opening Recommendation System
Another project I have worked on on-and-off. When I wanted to deliberately learn chess openings, I had to:

1. Search for openings that fit my "playstyle" (personally like playing aggressive positional openings), often through many Google searches
2. Ask higher rated players for input and advice
3. Learn the basics of the openings, mostly the most common lines
4. Play the opening in a few games and gauge how well I like the positions
5. Learn more of the sidelines of the openings
6. Play the opening in even more games and gauge how well I perform

There are too many chess openings and not enough time in the world for this process to be either effective or efficient. The clear path ahead was to make a recommender using [collaborative filtering](https://arxiv.org/pdf/2206.14312), maybe finding a way to add content-based filtering later on. So far I've only created a dataset using games from Chess.com to make user matrices, so the project is on hold.

- [Dataset](https://huggingface.co/datasets/juliusbroomfield/user-opening-interaction-matrix)
- [Code](https://github.com/juliusbroomfield/opening-recommendation-system) 

---

### 💻 Youtube Channel Exploratory Data Analysis
An EDA I worked on for an introductory programming class my freshman year. At the time I had a basic understanding of ML algorithms but hadn’t yet learned or worked with deep learning. We had the option to build a regression model, so I experimented with creating a neural network regression model. It performed pretty poorly (with an MSE around 45), but that was mostly due to limitations in the dataset.

- [Notebook](https://colab.research.google.com/drive/1ilrc7y1Z11Wqq-fNAYZWiqMgkMDaqxtJ?usp=sharing)
- [Code](https://github.com/juliusbroomfield/channel-growth-model/blob/main/analysis/youtube_models.ipynb)

---

