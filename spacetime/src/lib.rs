use spacetimedb::{spacetimedb, Identity, ReducerContext, Timestamp};

#[spacetimedb(table)]
pub struct User {
    #[primarykey]
    identity: Identity,
    name: Option<String>,
    online: bool,
}

#[spacetimedb(table)]
pub struct Message {
    sender: Identity,
    sent: Timestamp,
    text: String,
}

fn validate_name(name: String) -> Result<String, String> {
    if name.is_empty() {
        Err("Names must not be empty".to_string())
    } else {
        Ok(name)
    }
}

fn validate_message(name: String) -> Result<String, String> {
    if name.is_empty() {
        Err("Messages must not be empty".to_string())
    } else {
        Ok(name)
    }
}

#[spacetimedb(reducer)]
pub fn set_name(ctx: ReducerContext, name: String) -> Result<(), String> {
    let name = validate_name(name)?;
    if let Some(user) = User::filter_by_identity(&ctx.sender) {
        User::update_by_identity(
            &ctx.sender,
            User {
                name: Some(name),
                ..user
            },
        );
        Ok(())
    } else {
        Err("Cannot set name for unknown user".to_string())
    }
}

#[spacetimedb(reducer)]
pub fn send_message(ctx: ReducerContext, text: String) -> Result<(), String> {
    let text = validate_message(text)?;
    log::info!("{}", text);
    Message::insert(Message {
        sender: ctx.sender,
        sent: ctx.timestamp,
        text,
    });
    Ok(())
}

#[spacetimedb(connect)]
pub fn identity_connected(ctx: ReducerContext) {
    if let Some(user) = User::filter_by_identity(&ctx.sender) {
        User::update_by_identity(
            &ctx.sender,
            User {
                online: true,
                ..user
            },
        );
    } else {
        User::insert(User {
            online: true,
            identity: ctx.sender,
            name: None,
        })
        .unwrap();
    }
}

#[spacetimedb(disconnect)]
pub fn identity_disconnect(ctx: ReducerContext) {
    if let Some(user) = User::filter_by_identity(&ctx.sender) {
        User::update_by_identity(
            &ctx.sender,
            User {
                online: false,
                ..user
            },
        );
    } else {
        log::warn!("disconnect called for a user that doesn't yet exits");
    }
}
