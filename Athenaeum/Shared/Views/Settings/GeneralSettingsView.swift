/**
 GeneralSettingsView.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import SwiftUI

struct GeneralSettingsView: View {
    @AppStorage(HOST_DEFAULTS_KEY) private var host = ""
    @AppStorage(USERNAME_DEFAULTS_KEY) private var username = ""
    @AppStorage(PASSWORD_DEFAULTS_KEY) private var password = ""

    var body: some View {
        Form {
            HStack {
                Text("Athenaeum Server Host")
                Spacer()
                TextField("Host", text: $host)
            }
            HStack {
                Text("Athenaeum Server Username")
                Spacer()
                TextField("Username", text: $username)
            }
            HStack {
                Text("Athenaeum Server Password")
                Spacer()
                TextField("Password", text: $password)
            }
        }
        .padding(20)
        .frame(width: 450, height: 100)
    }
}
