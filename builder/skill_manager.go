package builder

import (
	// "ai_agents/vision_builder/utils"
	"context"
	"fmt"
	"strings"
	"sync"

	// 	"context"
	"io/ioutil"

	// "ai_agents/vision_builder/prompts"
	"ai_agents/vision_builder/utils"

	// "fmt"
	"log"
	"os"
	"path/filepath"

	// 	"strings"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms/openai"
	// "github.com/tmc/langchaingo/schema"

	// "github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type SkillManager struct {
	LLM               *openai.LLM
	Skills            map[string]map[string]string
	ControlPrimitives []string
	RetrievalTopK     int32
	CkptDir           string
	Vectordb          chroma.Store
	mu            sync.RWMutex // To safely access skills map from multiple goroutines
}

func NewSkillManager(modelName string, temperature float64, retrievalTopK int32, requestTimeout int, ckptDir string, resume bool) *SkillManager {
	store, errNs := chroma.New(
		chroma.WithChromaURL("http://localhost:8000"),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(uuid.New().String()),
	)
	if errNs != nil {
		log.Fatalf("new: %v\n", errNs)
	}

	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}
	utils.Mkdir(filepath.Join(ckptDir, "skill", "code"))
	utils.Mkdir(filepath.Join(ckptDir, "skill", "description"))
	utils.Mkdir(filepath.Join(ckptDir, "skill", "vectordb"))

	var skills map[string]map[string]string
	if resume {
		fmt.Printf("\033[33mLoading Skill Manager from %s/skill\033[0m\n", ckptDir)
		sklls, err := utils.LoadJSON(filepath.Join(ckptDir, "skill", "skills.json"))
		if err != nil {
			log.Fatal(err)
		}
		skills = sklls.(map[string]map[string]string)
	} else {
		skills = make(map[string]map[string]string)
	}

	controlPrimitives, err := LoadControlPrimitives(ckptDir,[]string{})
	if err != nil {
		log.Fatal(err)
	}
	  
	vectordb := store
	
	return &SkillManager{
		LLM:               llm,
		Skills:            skills,
		ControlPrimitives: controlPrimitives,
		RetrievalTopK:     retrievalTopK,
		CkptDir:           ckptDir,
		Vectordb:          vectordb,
	}
}

func LoadControlPrimitives(ckptPath string, primitiveNames []string) ([]string, error) {
	if primitiveNames == nil {
		primitiveNames = []string{}
		files, err := ioutil.ReadDir(filepath.Join(ckptPath, "control_primitives"))
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if filepath.Ext(file.Name()) == ".js" {
				primitiveNames = append(primitiveNames, file.Name()[:len(file.Name())-3])
			}
		}
	}

	var primitives []string
	for _, primitiveName := range primitiveNames {
		primitive, err := loadText(filepath.Join(ckptPath, "control_primitives", primitiveName+".go"))
		if err != nil {
			return nil, err
		}
		primitives = append(primitives, primitive)
	}
	return primitives, nil
}

func loadText(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}




// func (sm *SkillManager) Programs() string {
// 	var programs string
// 	for _, entry := range sm.Skills {
// 		programs += fmt.Sprintf("%s\n\n", entry["code"])
// 	}
// 	for _, primitives := range sm.ControlPrimitives {
// 		programs += fmt.Sprintf("%s\n\n", primitives)
// 	}
// 	return programs
// }

// func (sm *SkillManager) AddNewSkill(info map[string]interface{}) {
// 	if task, ok := info["task"].(string); ok && strings.HasPrefix(task, "Deposit useless items into the chest at") {
// 		// No need to reuse the deposit skill
// 		return
// 	}
// 	programName := info["program_name"].(string)
// 	programCode := info["program_code"].(string)
// 	dumpedProgramName := ""
// 	skillDescription := sm.GenerateSkillDescription(programName, programCode)
// 	fmt.Printf("\033[33mSkill Manager generated description for %s:\n%s\033[0m\n", programName, skillDescription)
// 	if _, exists := sm.Skills[programName]; exists {
// 		fmt.Printf("\033[33mSkill %s already exists. Rewriting!\033[0m\n", programName)
// 		sm.Vectordb.RemoveCollection()
// 		i := 2
// 		fNames, err := utils.ListDir(filepath.Join(sm.CkptDir, "skill", "code"))
// 		if err != nil{
// 			log.Fatalln(err)
// 		}
// 		for _, file := range fNames {
// 			if strings.HasPrefix(file.Name(), fmt.Sprintf("%sV", programName)) {
// 				i++
// 			}
// 		}
// 		dumpedProgramName = fmt.Sprintf("%sV%d", programName, i)
// 	} else {
// 		dumpedProgramName = programName
// 	}
// 	sm.Vectordb.AddDocuments(context.Background(),[]schema.Document{
// 		{
// 		PageContent: skillDescription,
// 		Metadata: map[string]any{"name": programName},
// 		Score: 0, 
// 		},
// 	})
		
// 	sm.Skills[programName] = map[string]string{
// 		"code":        programCode,
// 		"description": skillDescription,
// 	}

// 	colls, err := sm.Vectordb.SimilaritySearch(context.TODO(), programName,)
// 	if err != nil{
// 		log.Fatalln(err)
// 	}
// 	if len(colls) != len(sm.Skills){
// 		fmt.Println("vectordb is not synced with skills.json")
// 	}

// 	utils.DumpText(programCode, filepath.Join(sm.CkptDir, "skill", "code", fmt.Sprintf("%s.js", dumpedProgramName)))
// 	utils.DumpText(skillDescription, filepath.Join(sm.CkptDir, "skill", "description", fmt.Sprintf("%s.txt", dumpedProgramName)))
// 	utils.DumpJSON(sm.Skills, filepath.Join(sm.CkptDir, "skill", "skills.json"))
// }

// func (sm *SkillManager) GenerateSkillDescription(programName, programCode string) string {
// 	_, errAd := store.AddDocuments(context.Background(), []schema.Document{
// 		{PageContent: programName, Metadata: map[string]interface{}{programName:programCode}},
// 	})
	
	
// 	messages := []schema.ChatMessage{
// 		{Content: prompts.LoadPrompt(sm."skill")},
// 		{Content: programCode + "\n\n" + fmt.Sprintf("The main function is `%s`.", programName)},
// 	}
// 	skillDescription := fmt.Sprintf("    // %s", sm.LLM(messages))
// 	return fmt.Sprintf("async function %s(bot) {\n%s\n}", programName, skillDescription)
// }

// RetrieveSkills searches for skills similar to the query and returns their codes.
func (sm *SkillManager) RetrieveSkills(query string) []string {
    sm.mu.RLock() // Read lock for concurrent read safety
    defer sm.mu.RUnlock()
    if sm.RetrievalTopK == 0 {
        return []string{}
    }
    fmt.Printf("\033[33mSkill Manager retrieving for %d skills\033[0m\n", sm.RetrievalTopK)
	ctx := context.TODO()
    docsAndScores, err := sm.Vectordb.SimilaritySearch(ctx, query, int(sm.RetrievalTopK))
	if err != nil {
		log.Fatalln("Skill manager trying to retrieve similar skill from DB",err)
	}
    var skills []string
    var skillNames []string

    for _, docAndScore := range docsAndScores {
		codeMata := docAndScore.Metadata["name"].(map[string]any)
        skill, exists := codeMata["code"]
        if exists {
            skills = append(skills, skill.(string))
            skillNames = append(skillNames, docAndScore.PageContent)
        }
    }
    fmt.Printf("\033[33mSkill Manager retrieved skills: %s\033[0m\n", strings.Join(skillNames, ", "))

    return skills
}

func min(a, b int32) int32 {
    if a < b {
        return a
    }
    return b
}


