/**
 podcast.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Plot

func generateFeed() -> String {
    let location = "/Users/ckerson/samplebooks/"

    var feed = PodcastFeed(.title("Podcast Title"))

    return feed.render(indentedBy: .spaces(4))
}
