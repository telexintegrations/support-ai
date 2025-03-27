package chromadb

import (
	"context"
	"fmt"
	"strings"
)

func (c *ChromaDB) DeleteEntireOrganisationContext(ctx context.Context, orgID string) error {
	if orgID == "" {
		return ErrNoOrgId
	}

	collectionAsOrgId := fmt.Sprintf("%s-%s", collectionPrefix, orgID)

	col, err := c.ChromaDB().GetCollection(ctx, collectionAsOrgId, nil)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), fmt.Sprintf("%s does not exist", collectionAsOrgId)) {
			return ErrNoDataInOrg
		}
		return err
	}

	where := map[string]interface{}{
		"OrgId": orgID,
	}

	queryResults, err := col.Get(ctx, nil, where, nil, nil)
	if err != nil {
		fmt.Println("Error retrieving documents:", err)
		return err
	}

	if len(queryResults.Ids) == 0 {
		fmt.Printf("No documents found for OrgId: %s\n", orgID)
		return fmt.Errorf("No documents found for OrgId %s\n", orgID)
	}

	// Step 2: Delete by IDs
	doc, err := col.Delete(ctx, queryResults.Ids, where, nil)
	if err != nil {
		fmt.Println("Error deleting documents:", err)
		return err
	}

	fmt.Println("Doc ----------> ", doc)
	fmt.Printf("Successfully deleted all documents for OrgId: %s\n", orgID)
	return nil
}
