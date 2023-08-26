package appligen

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Applicant struct {
	FirstName      string        `yaml:"firstname"`
	LastName       string        `yaml:"lastname"`
	Street         string        `yaml:"street"`
	City           string        `yaml:"city"`
	Zip            string        `yaml:"zip"`
	Email          string        `yaml:"email"`
	Phone          string        `yaml:"phone"`
	PhoneFormatted string        `yaml:"phoneformatted"`
	LinkedIn       string        `yaml:"linkedin"`
	Applications   []Application `yaml:"applications"`
}

type Application struct {
	Company  string `yaml:"company"`
	Street   string `yaml:"street"`
	City     string `yaml:"city"`
	Country  string `yaml:"country"`
	Position string `yaml:"position"`
	Text     string `yaml:"text"`
	Zip      string `yaml:"zip"`
}

func copyFolder(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			err := os.Mkdir(destPath, 0755)
			if err != nil {
				return err
			}

			err = copyFolder(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			srcFile, err := os.Open(srcPath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewFromFile(templatePath string) (a Applicant, err error) {
	f, err := os.ReadFile(templatePath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(f, &a)

	return
}

func (a *Applicant) Generate(templatePath string) error {
	templateFile, err := os.Open(templatePath)
	if err != nil {
		return err
	}
	defer templateFile.Close()

	var templateContent strings.Builder
	_, err = io.Copy(&templateContent, templateFile)
	if err != nil {
		return err
	}

	for _, app := range a.Applications {
		folderName := fmt.Sprintf("%s_%s", app.Company, strings.ReplaceAll(app.Position, " ", "-"))

		if _, err := os.Stat(folderName); !os.IsNotExist(err) {
			err := os.RemoveAll(folderName)
			if err != nil {
				return err
			}
		}
		err := os.Mkdir(folderName, 0755)
		if err != nil {
			return err
		}

		templatePathDest := filepath.Join(folderName, "template.tex")
		templateFile, err := os.Create(templatePathDest)
		if err != nil {
			return err
		}
		defer templateFile.Close()

		_, err = io.WriteString(templateFile, templateContent.String())
		if err != nil {
			return err
		}

		dataPath := filepath.Join(folderName, "data.tex")
		dataTemplateContent := fmt.Sprintf(`\newcommand{\name}{%s %s}
\newcommand{\namerev}{%s, %s}
\newcommand{\sender}{%s %s, %s, %s %s}
\newcommand{\mail}{%s}
\newcommand{\linkedin}{%s}
\newcommand{\phone}{%s}
\newcommand{\phoneformatted}{%s}
\newcommand{\address}{
	%s

	%s

	%s %s

	%s
}
\newcommand{\position}{\emph{%s}}
\newcommand{\text}{
	%s
	\bigskip

}
`,
			a.FirstName,
			a.LastName,
			a.LastName,
			a.FirstName,
			a.FirstName,
			a.LastName,
			a.Street,
			a.Zip,
			a.City,
			a.Email,
			a.LinkedIn,
			a.Phone,
			a.PhoneFormatted,
			app.Company, app.Street, app.Zip, app.City, app.Country, app.Position, strings.ReplaceAll(app.Text, "\n\n", "\n\\bigskip\n\n"))

		dataFile, err := os.Create(dataPath)
		if err != nil {
			return err
		}
		defer dataFile.Close()

		_, err = io.WriteString(dataFile, dataTemplateContent)
		if err != nil {
			return err
		}

		err = copyFolder("figures", folderName)
		if err != nil {
			return err
		}

		err = runPDFLatex(folderName)
		if err != nil {
			return err
		}
	}
	return err
}

func runPDFLatex(folderName string) error {
	cmd := exec.Command("pdflatex", "template.tex")
	cmd.Dir = folderName
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pdflatex error: %v", err)
	}

	return nil
}
