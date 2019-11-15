"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// import Channel = require('./channel');
// import Message = require('./message');
var mongodb_1 = require("mongodb");
var channel_1 = require("./channel");
// import MongoClient = require('mongodb-typescript');
// import MongoObject = require('mongodb').ObjectID;
var assert = require('assert');
// Connection URL
var url = 'mongodb://localhost:27017';
// Database Name
var dbName = 'messaging';
// Create a new MongoClient
var client = new mongodb_1.MongoClient(url);
var db = null;
// openConnection does something
function openConnection() {
    // Use connect method to connect to the Server
    client.connect(function (err) {
        assert.equal(null, err);
        console.log("Connected successfully to server");
        db = client.db(dbName);
    });
    var general = new channel_1.Channel("general", "an open channel for all", false, [], "enter timestamp here", -1, "not yet edited");
    // channel that we always want at startup
    var result = insertNewChannel(general);
    if (result == null) {
        console.log("failed to create general channel upon opening connection to DB");
        // res.status(500);
    }
    return db;
}
db = openConnection();
// getAllChannels does something
function getAllChannels() {
    // if channels does not yet exist
    var cursor = db.channels.find();
    if (!cursor.hasNext()) {
        // Throw error
        console.log("No channels collection found");
        return null;
    }
    return cursor.forEach(function (m) { JSON.stringify(m); });
}
// insertNewChannel does something
function insertNewChannel(newChannel) {
    var result = db.channels.save({
        name: newChannel.name, description: newChannel.description,
        private: newChannel.private, members: newChannel.members,
        createdAt: newChannel.createdAt, creator: newChannel.creator,
        editedAt: newChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    newChannel._id = result._id;
    return newChannel;
}
// insertNewMessage does something
function insertNewMessage(newMessage) {
    if (newMessage.channelID == null) {
        return null;
    }
    var result = db.messages.save({
        channelID: newMessage.channelID, createdAt: newMessage.createdAt,
        body: newMessage.body, creator: newMessage.creator,
        editedAt: newMessage.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    newMessage._id = result._id;
    return newMessage;
}
// updatedChannel updates name and body of channel
function updatedChannel(existingChannel, req) {
    var result = db.channels.save({
        name: req.body.name, description: req.body.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    existingChannel.name = result.name;
    existingChannel.description = result.description;
    return existingChannel;
}
function addChannelMembers(existingChannel, req) {
    existingChannel.members.push(req.body.message.id);
    var result = db.channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    // Add the specified member
    // existingChannel.members = newMembers;
    return existingChannel;
}
function removeChannelMembers(existingChannel, req) {
    // Remove the specified member from this channel's list of members
    existingChannel.members.splice(req.body.message.id, 1);
    var result = db.channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    return existingChannel;
}
function updateMessage(existingMessage, req) {
    var result = db.messages.save({
        body: req.body, creator: existingMessage.creator,
        createdAt: existingMessage.createdAt, channelID: existingMessage.channelID,
        editedAt: existingMessage.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    existingMessage.body = result.body;
    return existingMessage;
}
// deleteChannel does something
function deleteChannel(existingChannel) {
    // We are not allowed to delete the general channel
    if (existingChannel.creator == -1) {
        return null;
    }
    db.channels.remove({ _id: new mongodb_1.ObjectID(existingChannel._id) });
    var result = db.messages.remove({ channelID: existingChannel._id });
    if (result.hasWriteError()) {
        return null;
    }
    return result;
}
function deleteMessage(existingMessage) {
    var result = db.messages.remove({ messageID: existingMessage._id });
    if (result.hasWriteError()) {
        return null;
    }
    return result;
}
// queryByChannelID does something
function getChannelByID(id) {
    if (id == null) {
        return null;
    }
    return db.channels.find({ _id: id });
}
function getMessageByID(id) {
    if (id == null) {
        return null;
    }
    return db.messages.find({ _id: id });
}
// last100Messages does something
function last100Messages(id) {
    if (id == null) {
        return null;
    }
    id = id.toString();
    return db.messages.find({ channelID: id }).sort({ createdAt: -1 }).limit(100);
}
// closeConnection does something
function closeConnection() {
    client.close();
}
//export the public functions
module.exports = {
    openConnection: openConnection,
    getAllChannels: getAllChannels,
    insertNewChannel: insertNewChannel,
    insertNewMessage: insertNewMessage,
    updatedChannel: updatedChannel,
    addChannelMembers: addChannelMembers,
    removeChannelMembers: removeChannelMembers,
    updateMessage: updateMessage,
    deleteChannel: deleteChannel,
    deleteMessage: deleteMessage,
    getChannelByID: getChannelByID,
    getMessageByID: getMessageByID,
    last100Messages: last100Messages,
    closeConnection: closeConnection
};
