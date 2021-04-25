/**
 BookDetailView.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import SwiftUI

struct BookDetailView: View {
    @ObservedObject var viewModel: BookDetailViewModel

    @ViewBuilder
    var body: some View {
        if viewModel.loading {
            ProgressView()
                .onAppear(perform: { viewModel.reload() })
        } else if let error = viewModel.error {
            Label(error.description, systemImage: "exclamationmark.triangle")
        } else if let item = viewModel.book {
            VStack(alignment: .leading) {
                Text(item.title)
                    .font(.title)
                    .lineLimit(3)
                Text("By \(item.authorString)")
                    .font(.title2)

                if let releaseDate = item.releaseDate {
                    Text("Released \(releaseDate)")
                        .font(.caption)
                }

                if let summary = item.summary {
                    Text(summary)
                }
            }
            .padding()
            Spacer()
        }
    }
}

//
// struct BookDetailView_Previews: PreviewProvider {
//    static var previews: some View {
//        BookDetailView()
//    }
// }
