package account

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/allincart-org/allincart-cli/extension"
	account_api "github.com/allincart-org/allincart-cli/internal/account-api"
	"github.com/allincart-org/allincart-cli/logging"
)

var accountCompanyProducerExtensionInfoPullCmd = &cobra.Command{
	Use:   "pull [path]",
	Short: "Generates local store configuration from account data",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		absolutePath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}

		zipExt, err := extension.GetExtensionByFolder(absolutePath)
		if err != nil {
			return fmt.Errorf("cannot open extension: %w", err)
		}

		zipName, err := zipExt.GetName()
		if err != nil {
			return fmt.Errorf("cannot get extension name: %w", err)
		}

		p, err := services.AccountClient.Producer(cmd.Context())
		if err != nil {
			return fmt.Errorf("cannot get producer endpoint: %w", err)
		}

		storeExt, err := p.GetExtensionByName(cmd.Context(), zipName)
		if err != nil {
			return fmt.Errorf("cannot get store extension: %w", err)
		}

		resourcesFolder := path.Join(zipExt.GetPath(), "src/Resources/store/")
		categoryList := make([]string, 0)
		localizations := make([]string, 0)
		tagsZH := make([]string, 0)
		tagsEN := make([]string, 0)
		videosZH := make([]string, 0)
		videosEN := make([]string, 0)
		highlightsZH := make([]string, 0)
		highlightsEN := make([]string, 0)
		featuresZH := make([]string, 0)
		featuresEN := make([]string, 0)
		faqZH := make([]extension.ConfigStoreFaq, 0)
		faqEN := make([]extension.ConfigStoreFaq, 0)

		if _, err := os.Stat(resourcesFolder); os.IsNotExist(err) {
			err = os.MkdirAll(resourcesFolder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("cannot create file: %w", err)
			}
		}

		var iconConfigPath *string

		if len(storeExt.IconURL) > 0 {
			icon := "src/Resources/store/icon.png"
			iconConfigPath = &icon
			err := downloadFileTo(cmd.Context(), storeExt.IconURL, path.Join(resourcesFolder, "icon.png"))
			if err != nil {
				return fmt.Errorf("cannot download file: %w", err)
			}
		}

		if storeExt.Category != nil {
			categoryList = append(categoryList, storeExt.Category.Name)
		} else {
			for _, category := range storeExt.Categories {
				categoryList = append(categoryList, category.Name)
				break
			}
		}

		for _, localization := range storeExt.Localizations {
			localizations = append(localizations, localization.Name)
		}

		storeImages, err := p.GetExtensionImages(cmd.Context(), storeExt.Id)
		if err != nil {
			return fmt.Errorf("cannot get extension images: %w", err)
		}

		if len(storeImages) > 0 {
			imagesDir := path.Join(zipExt.GetPath(), "src/Resources/store/images/")

			if err := writeImages(cmd.Context(), imagesDir, 0, storeImages); err != nil {
				return fmt.Errorf("cannot write images: %w", err)
			}

			if err := writeImages(cmd.Context(), imagesDir, 1, storeImages); err != nil {
				return fmt.Errorf("cannot write images: %w", err)
			}
		}

		chineseDescription := ""
		englishDescription := ""
		chineseInstallationManual := ""
		englishInstallationManual := ""
		chineseMetaTitle := ""
		englishMetaTitle := ""
		chineseMetaDescription := ""
		englishMetaDescription := ""

		for _, info := range storeExt.Infos {
			language := info.Locale.Name[0:2]

			if language == "zh" {
				chineseDescription = "file:src/Resources/store/description.zh.html"
				chineseInstallationManual = "file:src/Resources/store/installation_manual.zh.html"
				chineseMetaTitle = info.MetaTitle
				chineseMetaDescription = info.MetaDescription

				if err := os.WriteFile(path.Join(zipExt.GetPath(), chineseDescription[5:]), []byte(info.Description), os.ModePerm); err != nil {
					return fmt.Errorf("cannot write file: %w", err)
				}

				if err := os.WriteFile(path.Join(zipExt.GetPath(), chineseInstallationManual[5:]), []byte(info.InstallationManual), os.ModePerm); err != nil {
					return fmt.Errorf("cannot write file: %w", err)
				}

				for _, element := range info.Tags {
					tagsZH = append(tagsZH, element.Name)
				}

				for _, element := range info.Videos {
					videosZH = append(videosZH, element.URL)
				}

				if info.Highlights != "" {
					highlightsZH = append(highlightsZH, strings.Split(info.Highlights, "\n")...)
				}
				if info.Features != "" {
					featuresZH = append(featuresZH, strings.Split(info.Features, "\n")...)
				}

				for _, element := range info.Faqs {
					faqZH = append(faqZH, extension.ConfigStoreFaq{Question: element.Question, Answer: element.Answer})
				}
			} else {
				englishDescription = "file:src/Resources/store/description.en.html"
				englishInstallationManual = "file:src/Resources/store/installation_manual.en.html"
				englishMetaTitle = info.MetaTitle
				englishMetaDescription = info.MetaDescription

				if err := os.WriteFile(path.Join(zipExt.GetPath(), englishDescription[5:]), []byte(info.Description), os.ModePerm); err != nil {
					return fmt.Errorf("cannot write file: %w", err)
				}

				if err := os.WriteFile(path.Join(zipExt.GetPath(), englishInstallationManual[5:]), []byte(info.InstallationManual), os.ModePerm); err != nil {
					return fmt.Errorf("cannot write file: %w", err)
				}

				for _, element := range info.Tags {
					tagsEN = append(tagsEN, element.Name)
				}

				for _, element := range info.Videos {
					videosEN = append(videosEN, element.URL)
				}

				if info.Highlights != "" {
					highlightsEN = append(highlightsEN, strings.Split(info.Highlights, "\n")...)
				}

				if info.Features != "" {
					featuresEN = append(featuresEN, strings.Split(info.Features, "\n")...)
				}

				for _, element := range info.Faqs {
					faqEN = append(faqEN, extension.ConfigStoreFaq{Question: element.Question, Answer: element.Answer})
				}
			}
		}

		extType := "extension"

		if storeExt.ProductType != nil {
			extType = storeExt.ProductType.Name
		}

		newCfg := zipExt.GetExtensionConfig()

		newCfg.Store.Icon = iconConfigPath
		newCfg.Store.DefaultLocale = &storeExt.StandardLocale.Name
		newCfg.Store.Type = &extType
		newCfg.Store.AutomaticBugfixVersionCompatibility = &storeExt.AutomaticBugfixVersionCompatibility
		newCfg.Store.Localizations = &localizations
		newCfg.Store.Description = extension.ConfigTranslated[string]{Chinese: &chineseDescription, English: &englishDescription}
		newCfg.Store.InstallationManual = extension.ConfigTranslated[string]{Chinese: &chineseInstallationManual, English: &englishInstallationManual}
		newCfg.Store.Categories = &categoryList
		newCfg.Store.Tags = extension.ConfigTranslated[[]string]{Chinese: &tagsZH, English: &tagsEN}
		newCfg.Store.Videos = extension.ConfigTranslated[[]string]{Chinese: &videosZH, English: &videosEN}
		newCfg.Store.Highlights = extension.ConfigTranslated[[]string]{Chinese: &highlightsZH, English: &highlightsEN}
		newCfg.Store.Features = extension.ConfigTranslated[[]string]{Chinese: &featuresZH, English: &featuresEN}
		newCfg.Store.Faq = extension.ConfigTranslated[[]extension.ConfigStoreFaq]{Chinese: &faqZH, English: &faqEN}
		newCfg.Store.MetaTitle = extension.ConfigTranslated[string]{Chinese: &chineseMetaTitle, English: &englishMetaTitle}
		newCfg.Store.MetaDescription = extension.ConfigTranslated[string]{Chinese: &chineseMetaDescription, English: &englishMetaDescription}
		newCfg.Store.Images = nil

		if len(storeImages) > 0 {
			imageDir := "src/Resources/store/images"
			newCfg.Store.ImageDirectory = &imageDir
		}

		content, err := yaml.Marshal(newCfg)
		if err != nil {
			return fmt.Errorf("cannot encode yaml: %w", err)
		}

		extCfgFile := fmt.Sprintf("%s/%s", zipExt.GetPath(), newCfg.FileName)
		err = os.WriteFile(extCfgFile, content, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot save file: %w", err)
		}

		logging.FromContext(cmd.Context()).Infof("Files has been written to the given extension folder")

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionInfoCmd.AddCommand(accountCompanyProducerExtensionInfoPullCmd)
}

func downloadFileTo(ctx context.Context, url string, target string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logging.FromContext(ctx).Errorf("downloadFileTo: %v", err)
		}
	}()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read file body: %w", err)
	}

	err = os.WriteFile(target, content, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	return nil
}

func writeImages(ctx context.Context, imagePath string, index int, storeImages []*account_api.ExtensionImage) error {
	imageMap := make(map[int]string)

	for _, image := range storeImages {
		if image.Details[index].Activated {
			priority := image.Priority

			if _, ok := imageMap[priority]; !ok {
				imageMap[priority] = image.RemoteLink
			} else {
				for {
					priority++
					if _, ok := imageMap[priority]; !ok {
						imageMap[priority] = image.RemoteLink
						break
					}
				}
			}
		}
	}

	if index == 0 {
		imagePath = path.Join(imagePath, "zh")
	} else {
		imagePath = path.Join(imagePath, "en")
	}

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(imagePath, os.ModePerm); err != nil {
			return err
		}
	}

	for priority, link := range imageMap {
		if err := downloadFileTo(ctx, link, path.Join(imagePath, fmt.Sprintf("%d.png", priority))); err != nil {
			return err
		}
	}

	return nil
}
