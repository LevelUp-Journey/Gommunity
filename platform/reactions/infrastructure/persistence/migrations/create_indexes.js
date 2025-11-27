// MongoDB indexes for reactions collection
// Run this in MongoDB shell or using mongosh

db.reactions.createIndex(
  { post_id: 1, user_id: 1 },
  { 
    unique: true,
    name: "idx_unique_post_user_reaction",
    background: true
  }
);

db.reactions.createIndex(
  { post_id: 1, created_at: -1 },
  { 
    name: "idx_post_reactions_timeline",
    background: true
  }
);

db.reactions.createIndex(
  { post_id: 1, reaction_type: 1 },
  { 
    name: "idx_post_reaction_type",
    background: true
  }
);

db.reactions.createIndex(
  { user_id: 1, created_at: -1 },
  { 
    name: "idx_user_reactions_timeline",
    background: true
  }
);
