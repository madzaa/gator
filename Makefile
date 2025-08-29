.PHONY: build run reset seed

# Build the application
build:
	go build -o gator

# Run the application
run:
	go run .

# Reset database (delete all data)
reset:
	go run . reset

# Seed sample data
seed: reset
	@( \
		go run . register alice && \
		go run . addfeed "BBC News" "http://feeds.bbci.co.uk/news/rss.xml" && \
		go run . addfeed "NYT Home Page" "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml" && \
		go run . register bob && \
		go run . follow "http://feeds.bbci.co.uk/news/rss.xml" && \
		go run . addfeed "Hacker News" "https://news.ycombinator.com/rss" && \
		go run . follow "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml" && \
		go run . register carol && \
		go run . addfeed "Reddit Programming" "https://www.reddit.com/r/programming/.rss" && \
		go run . follow "https://news.ycombinator.com/rss" && \
		go run . follow "http://feeds.bbci.co.uk/news/rss.xml" && \
		go run . follow "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml" && \
		go run . follow "https://www.reddit.com/r/programming/.rss" \
	)