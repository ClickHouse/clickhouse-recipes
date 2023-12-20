FROM logs 
SELECT count(), 
       formatReadableQuantity(count()) AS countFriendly, 
       now() 
Format PrettyNoEscapes;
