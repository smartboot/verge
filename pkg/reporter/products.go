package reporter

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ibuilding-x/driver-box/pkg/driverbox/helper"
	"go.uber.org/zap"

	"github.com/smartboot/verge-export/pkg/rpc"
)

func (r *Reporter) ReportProducts(products []rpc.ProductInfo) error {
	helper.Logger.Info("reporting products", zap.Int("productCount", len(products)))
	return r.postReport("report/products", products)
}

// parseProductID 从文件名解析产品ID和模型ID
func parseProductID(filename string) (productID, modelID string) {
	// 移除文件扩展名
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 按冒号分割
	parts := strings.Split(name, ":")
	if len(parts) == 2 {
		productID = parts[0]
		modelID = parts[1]
	}
	return
}

// calculateMD5 计算文件的MD5值
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CollectAndReportProducts collects product information from library and reports to server
func (r *Reporter) CollectAndReportProducts() error {
	// Collect and process model and driver files to generate product list
	productMap := make(map[string]*rpc.ProductInfo) // productID -> ProductInfo

	// Process model files
	modelDir := "res/library/model"
	if err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			filename := filepath.Base(path)
			productID, modelID := parseProductID(filename)
			if productID == "" {
				return nil // skip files that can't be parsed
			}

			md5, err := calculateMD5(path)
			if err != nil {
				helper.Logger.Error("Failed to calculate MD5 for model file", zap.String("path", path), zap.Error(err))
				return nil
			}

			// Initialize product info if not exists
			if productMap[productID] == nil {
				productMap[productID] = &rpc.ProductInfo{
					Product: productID,
					Models:  make(map[string]string),
					Driver:  make(map[string]string),
				}
			}

			// Add model
			productMap[productID].Models[modelID] = md5
		}
		return nil
	}); err != nil {
		helper.Logger.Error("Failed to process model files", zap.Error(err))
		return err
	}

	// Process driver files
	driverDir := "res/library/driver"
	if err := filepath.Walk(driverDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".lua" {
			filename := filepath.Base(path)
			productID, driverID := parseProductID(filename)
			if productID == "" {
				return nil // skip files that can't be parsed
			}

			md5, err := calculateMD5(path)
			if err != nil {
				helper.Logger.Error("Failed to calculate MD5 for driver file", zap.String("path", path), zap.Error(err))
				return nil
			}

			// Initialize product info if not exists
			if productMap[productID] == nil {
				productMap[productID] = &rpc.ProductInfo{
					Product: productID,
					Models:  make(map[string]string),
					Driver:  make(map[string]string),
				}
			}

			// Add driver
			productMap[productID].Driver[driverID] = md5
		}
		return nil
	}); err != nil {
		helper.Logger.Error("Failed to process driver files", zap.Error(err))
		return err
	}

	// Generate product list and calculate final hash
	products := make([]rpc.ProductInfo, 0)
	for _, productInfo := range productMap {
		// Collect all model and driver hashes for final hash calculation
		var allHashes []string

		// Collect all model hashes
		for _, hash := range productInfo.Models {
			allHashes = append(allHashes, hash)
		}

		// Collect all driver hashes
		for _, hash := range productInfo.Driver {
			allHashes = append(allHashes, hash)
		}

		// Sort all hashes for consistent final hash calculation
		sort.Strings(allHashes)

		// Concatenate sorted hashes and calculate final hash
		var concatenatedHashes strings.Builder
		for _, hash := range allHashes {
			concatenatedHashes.WriteString(hash)
		}
		productInfo.Hash = fmt.Sprintf("%x", md5.Sum([]byte(concatenatedHashes.String())))

		products = append(products, *productInfo)
	}

	// Report products
	return r.ReportProducts(products)
}
