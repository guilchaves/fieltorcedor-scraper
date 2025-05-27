package notifiedgames

import (
    "bufio"
    "os"
    "sync"
)

type FileNotifiedGamesRepository struct {
    notifiedGames map[string]bool
    mu            sync.RWMutex
    filePath      string
}

func NewFileNotifiedGamesRepository(filePath string) *FileNotifiedGamesRepository {
    repo := &FileNotifiedGamesRepository{
        notifiedGames: make(map[string]bool),
        filePath:      filePath,
    }
    repo.loadFromFile()
    return repo
}

func (r *FileNotifiedGamesRepository) SaveNotifiedGame(gameID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.notifiedGames[gameID] = true
    return r.saveToFile(gameID)
}

func (r *FileNotifiedGamesRepository) IsGameNotified(gameID string) (bool, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    _, ok := r.notifiedGames[gameID]
    return ok, nil
}

func (r *FileNotifiedGamesRepository) loadFromFile() error {
    file, err := os.Open(r.filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        gameID := scanner.Text()
        r.notifiedGames[gameID] = true
    }
    return scanner.Err()
}

func (r *FileNotifiedGamesRepository) saveToFile(gameID string) error {
    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = file.WriteString(gameID + "\n")
    return err
}
